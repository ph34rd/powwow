package services

import (
	"time"

	"github.com/ph34rd/powwow/pkg/limiter"

	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"

	"github.com/ph34rd/powwow/pkg/complexer"
	"github.com/ph34rd/powwow/pkg/pow"
	"github.com/ph34rd/powwow/pkg/session/manager"
	"github.com/ph34rd/powwow/pkg/stat"
	"github.com/ph34rd/powwow/pkg/wow"
)

const defaultGracePeriod = 30 * time.Second

type Services struct {
	PoW       pow.PoW
	WoW       wow.WoW
	Sampler   stat.CPUSampler
	Complexer complexer.Complexer
	Manager   manager.Manager
	Limiter   limiter.StringLimiter
}

func NewServices(lg *zap.Logger) (*Services, error) {
	smSrv := stat.NewCPUSampler(lg, time.Second)
	services := &Services{
		PoW:       pow.New(sha3.New256),
		WoW:       wow.NewInMemory(),
		Sampler:   smSrv,
		Complexer: complexer.NewComplexer(smSrv),
		Manager:   manager.NewManager(),
		Limiter:   limiter.NewStringLimiter(),
	}
	return services, nil
}

func (s *Services) Shutdown() {
	s.Manager.Close()
	select {
	case <-time.After(defaultGracePeriod):
		s.Manager.Shutdown()
		<-s.Manager.Done()
	case <-s.Manager.Done():
	}
	s.Sampler.Shutdown()
}
