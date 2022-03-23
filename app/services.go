package app

import (
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/ph34rd/powwow/pkg/complexer"
	"github.com/ph34rd/powwow/pkg/limiter"
	"github.com/ph34rd/powwow/pkg/logger"
	"github.com/ph34rd/powwow/pkg/pow"
	"github.com/ph34rd/powwow/pkg/session/manager"
	"github.com/ph34rd/powwow/pkg/stat"
	"github.com/ph34rd/powwow/pkg/wow"
)

type Services struct {
	Logger    logger.Logger
	PoW       pow.PoW
	WoW       wow.WoW
	Sampler   stat.CPUSampler
	Complexer complexer.Complexer
	Manager   manager.Manager
	Limiter   limiter.StringLimiter
}

func NewServices(lg logger.Logger) (*Services, error) {
	smSrv := stat.NewCPUSampler(lg, time.Second)
	services := &Services{
		Logger:    lg,
		PoW:       pow.New(sha3.New256),
		WoW:       wow.NewInMemory(),
		Sampler:   smSrv,
		Complexer: complexer.NewComplexer(smSrv),
		Manager:   manager.NewManager(),
		Limiter:   limiter.NewStringLimiter(),
	}
	return services, nil
}
