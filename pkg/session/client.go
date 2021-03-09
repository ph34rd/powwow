package session

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/ph34rd/powwow/pb"
	"github.com/ph34rd/powwow/pkg/pow/hashcash"
	"github.com/ph34rd/powwow/pkg/transport"
)

const defaultPingDuration = 10 * time.Second

type ClientSession interface {
	GetWoW(ctx context.Context) (string, error)
}

// ClientHandler implements ClientSession.
type ClientHandler struct {
	conn      net.Conn
	r         *bufio.Reader
	w         *bufio.Writer
	transport transport.Transport

	challenge  []byte
	complexity uint8
	mintNonce  []byte
	mintErr    error
	wow        string
}

func NewClientHandler(conn net.Conn) *ClientHandler {
	timeoutConn := newConnDeadlineWrapper(conn, defaultTimeout)
	r := newBufioReader(timeoutConn)
	w := newBufioWriter(timeoutConn)
	return &ClientHandler{
		conn:      conn,
		r:         r,
		w:         w,
		transport: transport.NewTLVTransport(r, w, 0),
	}
}

func (c *ClientHandler) GetWoW(ctx context.Context) (string, error) {
	err := c.transport.NextReader(func(size uint32, typeID uint16, r io.Reader) error {
		if typeID == pbServerHandshake {
			return c.handleHandshake(r)
		}
		return fmt.Errorf("unknown message: %d", typeID)
	})
	if err != nil {
		return "", err
	}

	t := time.NewTicker(defaultPingDuration)
	mintCtx, mintCancel := context.WithCancel(ctx)
	mintDone := c.mint(mintCtx)
loop:
	for {
		select {
		case <-t.C:
			err = c.transport.Ping() // ping server while we mint
			if err != nil {
				mintCancel()
				return "", err
			}
			err = c.transport.NextReader(func(size uint32, typeID uint16, r io.Reader) error {
				return fmt.Errorf("unknown message: %d", typeID)
			})
			if err != nil {
				mintCancel()
				return "", err
			}
		case <-mintDone:
			mintCancel()
			break loop
		}
	}
	if c.mintErr != nil {
		return "", c.mintErr
	}

	err = c.sendNonce()
	if err != nil {
		return "", err
	}
	err = c.sendWoWReq()
	if err != nil {
		return "", err
	}

	// read all packets if any
	for {
		err = c.transport.NextReader(func(size uint32, typeID uint16, r io.Reader) error {
			if typeID == pbWoWResponse {
				return c.handleWoWResp(r)
			}
			return fmt.Errorf("unknown message: %d", typeID)
		})
		if err != nil {
			return "", err
		}
		if len(c.wow) > 0 {
			break
		}
	}

	c.transport.Close()
	return c.wow, nil
}

func (c *ClientHandler) handleHandshake(r io.Reader) error {
	var pkt pb.ServerHandshake
	err := unmarshalHelper(r, &pkt)
	if err != nil {
		return err
	}
	c.challenge = pkt.Challenge
	c.complexity = uint8(pkt.Complexity)
	return nil
}

func (c *ClientHandler) mint(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		iter, err := hashcash.NewFastIter(runtime.NumCPU())
		if err != nil {
			return
		}
		minter, err := hashcash.NewParallel(sha3.New256, iter, int(c.complexity), runtime.NumCPU())
		if err != nil {
			return
		}
		c.mintNonce, c.mintErr = minter.Mint(ctx, c.challenge)
		close(ch)
	}()
	return ch
}

func (c *ClientHandler) sendNonce() error {
	m := &pb.ClientHandshake{Nonce: c.mintNonce}
	typeID, err := pbMessageType(m)
	if err != nil {
		return err
	}
	err = c.transport.NextWriter(uint32(m.Size()), typeID, frameWriterFunc(m))
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientHandler) sendWoWReq() error {
	m := &pb.WoWRequest{}
	typeID, err := pbMessageType(m)
	if err != nil {
		return err
	}
	err = c.transport.NextWriter(uint32(m.Size()), typeID, frameWriterFunc(m))
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientHandler) handleWoWResp(r io.Reader) error {
	var pkt pb.WoWResponse
	err := unmarshalHelper(r, &pkt)
	if err != nil {
		return err
	}
	c.wow = pkt.Wow
	return nil
}
