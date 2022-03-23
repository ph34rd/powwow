package session

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/ph34rd/powwow/pb"
	"github.com/ph34rd/powwow/pkg/logger"
	"github.com/ph34rd/powwow/pkg/session/manager"
	"github.com/ph34rd/powwow/pkg/session/transport"
)

var errUnverifiedClient = errors.New("unverified client, access denied")

const defaultTimeout = 15 * time.Second

type ServerSession interface {
	Handle()
	Stop()
}

// ServerHandler implements ServerSession.
type ServerHandler struct {
	logger   logger.Logger
	conn     net.Conn
	services *ServerServices

	mu   sync.Mutex
	drop bool
	once sync.Once

	challenge  []byte
	complexity uint8
	verified   bool
	wowDone    bool
	limiter    *rate.Limiter
	r          *bufio.Reader
	w          *bufio.Writer
	transport  transport.Transport
}

func NewServerSession(conn net.Conn, services *ServerServices) *ServerHandler {
	sessLogger := services.Logger.With(
		zap.String("id", xid.New().String()),
		zap.String("remote", conn.RemoteAddr().String()),
		zap.String("local", conn.LocalAddr().String()),
	)
	timeoutConn := newConnDeadlineWrapper(conn, defaultTimeout)
	r := newBufioReader(timeoutConn)
	w := newBufioWriter(timeoutConn)
	return &ServerHandler{
		logger:    sessLogger,
		conn:      conn,
		services:  services,
		r:         r,
		w:         w,
		transport: transport.NewTLVTransport(r, w, 0),
	}
}

func (s *ServerHandler) Handle() {
	adr := manager.Stopper(s)
	adrPtr := &adr

	defer func() {
		s.close()
		putBufioReader(s.r)
		putBufioWriter(s.w)
		err := s.services.Manager.Track(adrPtr, false)
		if err != nil {
			s.logger.Error("track error", zap.Error(err))
		}
		s.limitSpend()
		s.logger.Info("client disconnected")
	}()
	s.logger.Info("client connected")

	err := s.services.Manager.Track(adrPtr, true)
	if err != nil {
		s.logger.Error("track error", zap.Error(err))
		s.setDrop(true)
		return
	}

	if s.limitReached() {
		s.logger.Info("connection limit reached")
		s.setDrop(true)
		return
	}

	err = s.sendHandshake()
	if err != nil {
		s.setDrop(true)
		return
	}

	for {
		err = s.transport.NextReader(func(size uint32, typeID uint16, r io.Reader) error {
			switch typeID {
			case pbClientHandshake:
				return s.verifyClient(r)
			case pbWoWRequest:
				return s.handleWoW(r)
			}
			return fmt.Errorf("unknown message: %d", typeID)
		})
		if err != nil {
			if err != transport.ErrClose {
				s.setDrop(true)
				s.logger.Debug("frame read error", zap.Error(err))
			}
			return
		}
		if s.wowDone {
			err = s.transport.Close()
			if err != nil && err != transport.ErrClose {
				s.setDrop(true)
			}
			return
		}
	}
}

func (s *ServerHandler) setDrop(v bool) {
	s.mu.Lock()
	s.drop = v
	s.mu.Unlock()
}

func (s *ServerHandler) close() {
	var drop bool
	s.mu.Lock()
	drop = s.drop
	s.mu.Unlock()

	s.once.Do(func() {
		if drop {
			if l, ok := s.conn.(linger); ok {
				l.SetLinger(0)
			}
		}
		s.conn.Close()
	})
}

func (s *ServerHandler) Stop() {
	s.setDrop(true)
	s.close()
}

func (s *ServerHandler) makeHandshake() (*pb.ServerHandshake, uint16, error) {
	var err error
	s.challenge, err = s.services.PoW.Challenge()
	if err != nil {
		s.logger.Error("create challenge error", zap.Error(err))
		return nil, 0, err
	}
	s.complexity = s.services.Complexer.NextComplexity()
	m := pb.ServerHandshake{Challenge: s.challenge, Complexity: uint32(s.complexity)}
	typeID, err := pbMessageType(&m)
	if err != nil {
		s.logger.Error("type error", zap.Error(err))
		return nil, 0, err
	}
	return &m, typeID, nil
}

func (s *ServerHandler) sendHandshake() error {
	pkt, typeID, err := s.makeHandshake()
	if err != nil {
		return err
	}
	err = s.transport.NextWriter(uint32(pkt.Size()), typeID, frameWriterFunc(pkt))
	if err != nil {
		s.logger.Debug("frame write error", zap.Error(err))
		return err
	}
	s.logger.Debug("frame write", zap.String("frame", pkt.String()))
	return nil
}

func (s *ServerHandler) verifyClient(r io.Reader) error {
	var pkt pb.ClientHandshake
	err := unmarshalHelper(r, &pkt)
	if err != nil {
		s.logger.Debug("frame read error", zap.Error(err))
		return err
	}
	s.logger.Debug("frame read", zap.String("frame", pkt.String()))

	err = s.services.PoW.Verify(s.challenge, pkt.Nonce, s.complexity)
	if err != nil {
		s.logger.Info("nonce verify error", zap.Error(err))
		return err
	}
	s.verified = true
	s.logger.Info("nonce verified", zap.ByteString("nonce", pkt.Nonce))
	return nil
}

func (s *ServerHandler) handleWoW(r io.Reader) error {
	if !s.verified {
		return errUnverifiedClient
	}
	var pkt pb.WoWRequest
	err := unmarshalHelper(r, &pkt)
	if err != nil {
		s.logger.Debug("frame read error", zap.Error(err))
		return err
	}
	s.logger.Debug("frame read", zap.String("frame", pkt.String()))

	wowResult, err := s.services.WoW.GetNext(context.Background())
	if err != nil {
		s.logger.Error("wow service error", zap.Error(err))
		return err
	}
	s.logger.Debug("wow service result", zap.String("wow", wowResult))

	m := pb.WoWResponse{Wow: wowResult}
	typeID, err := pbMessageType(&m)
	if err != nil {
		s.logger.Error("type error", zap.Error(err))
		return err
	}
	err = s.transport.NextWriter(uint32(m.Size()), typeID, frameWriterFunc(&m))
	if err != nil {
		s.logger.Debug("frame write error", zap.Error(err))
		return err
	}
	s.logger.Debug("frame write", zap.String("frame", m.String()))
	s.wowDone = true
	s.logger.Info("client wow complete", zap.String("wow", wowResult))
	return nil
}

func (s *ServerHandler) limitReached() bool {
	ip := extractIP(s.conn)
	if len(ip) == 0 {
		return false
	}
	s.limiter = s.services.Limiter.Get(ip)
	t := time.Now()
	reservation := s.limiter.ReserveN(t, 1)
	defer reservation.CancelAt(t)
	return reservation.Delay() > 0
}

func (s *ServerHandler) limitSpend() {
	if s.limiter != nil && !s.verified {
		s.limiter.Allow()
	}
}
