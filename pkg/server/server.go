package server

import (
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ph34rd/powwow/pkg/server/pool"
	"github.com/ph34rd/powwow/pkg/server/reuse"
	"github.com/ph34rd/powwow/pkg/services"
	"github.com/ph34rd/powwow/pkg/session"
)

const minWorkers = 65536

type Server struct {
	logger *zap.Logger

	inShutdown sync.Once
	closeCh    chan struct{}

	mu       sync.Mutex
	ln       net.Listener
	connList map[net.Conn]struct{}

	bind     string
	services *services.Services
}

func NewServer(lg *zap.Logger, bind string) (*Server, error) {
	srv, err := services.NewServices(lg)
	lg.Info("services initialized")
	if err != nil {
		lg.Error("services init error", zap.Error(err))
		return nil, err
	}
	return &Server{
		logger:   lg,
		closeCh:  make(chan struct{}),
		bind:     bind,
		services: srv,
	}, nil
}

func (s *Server) Run() error {
	ln, err := reuse.Listen("tcp", s.bind)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.ln = &onceCloseListener{Listener: ln}
	s.mu.Unlock()
	defer s.ln.Close()

	s.logger.Info("listen started", zap.String("bind", s.bind))
	defer func() {
		s.logger.Info("listen stopped")
		s.afterListenerClose()
	}()

	var tempDelay time.Duration // how long to sleep on accept failure
	wPool, err := pool.NewWorkerPool(minWorkers)
	if err != nil {
		s.logger.Info("worker pool error", zap.Error(err))
		return err
	}
	wPool.Init()

loop:
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.closeCh:
				break loop
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.logger.Error("accept error, retrying", zap.Error(err), zap.Duration("delay", tempDelay))
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		wPool.Go(session.NewServerSession(s.logger, s.services, conn).Handle)
	}
	return nil
}

func (s *Server) afterListenerClose() {
	s.services.Shutdown()
	s.logger.Info("all services stopped")
}

func (s *Server) stop() {
	close(s.closeCh)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ln != nil {
		s.ln.Close()
	}
}

func (s *Server) Stop() {
	s.inShutdown.Do(func() {
		s.logger.Info("server shutdown")
		s.stop()
	})
}
