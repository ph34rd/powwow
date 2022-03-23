package app

import (
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/ph34rd/powwow/pkg/complexer"
	"github.com/ph34rd/powwow/pkg/limiter"
	"github.com/ph34rd/powwow/pkg/logger"
	"github.com/ph34rd/powwow/pkg/pow"
	"github.com/ph34rd/powwow/pkg/session"
	"github.com/ph34rd/powwow/pkg/session/manager"
	"github.com/ph34rd/powwow/pkg/stat"
	"github.com/ph34rd/powwow/pkg/wow"
)

const defaultGracePeriod = 30 * time.Second

type lifecycle struct {
	*session.ServerServices
}

func newLifecycle(lg logger.Logger) (*lifecycle, error) {
	smSrv := stat.NewCPUSampler(lg, time.Second)
	services := &session.ServerServices{
		Logger:    lg,
		PoW:       pow.New(sha3.New256),
		WoW:       wow.NewInMemory(),
		Sampler:   smSrv,
		Complexer: complexer.NewComplexer(smSrv),
		Manager:   manager.NewManager(),
		Limiter:   limiter.NewStringLimiter(),
	}
	return &lifecycle{ServerServices: services}, nil
}

func (s *lifecycle) Shutdown() {
	s.Manager.Close()
	select {
	case <-time.After(defaultGracePeriod):
		s.Manager.Shutdown()
		<-s.Manager.Done()
	case <-s.Manager.Done():
	}
	s.Sampler.Shutdown()
}
