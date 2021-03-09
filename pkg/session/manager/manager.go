package manager

import (
	"errors"
	"sync"
)

var ErrManagerClosed = errors.New("manager: closed")

type Stopper interface {
	Stop()
}

type Manager interface {
	Track(*Stopper, bool) error
	Done() <-chan struct{}
	Close()
	Shutdown()
}

type Impl struct {
	mu         sync.Mutex
	activeConn map[*Stopper]struct{}
	inShutdown bool
	doneCh     chan struct{}
}

func NewManager() *Impl {
	return &Impl{activeConn: make(map[*Stopper]struct{}), doneCh: make(chan struct{})}
}

func (s *Impl) Track(stp *Stopper, add bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if add {
		if s.inShutdown {
			return ErrManagerClosed
		}
		s.activeConn[stp] = struct{}{}
	} else {
		delete(s.activeConn, stp)
	}
	if s.inShutdown && len(s.activeConn) == 0 {
		close(s.doneCh)
	}
	return nil
}

func (s *Impl) Close() {
	s.mu.Lock()
	s.closeLocked()
	s.mu.Unlock()
}

func (s *Impl) closeLocked() {
	if s.inShutdown {
		return
	}
	s.inShutdown = true
	if len(s.activeConn) == 0 {
		close(s.doneCh)
	}
}

func (s *Impl) Done() <-chan struct{} {
	return s.doneCh
}

func (s *Impl) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeLocked()
	for c := range s.activeConn {
		(*c).Stop()
	}
}
