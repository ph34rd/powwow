package session

import (
	"context"

	"github.com/ph34rd/powwow/pkg/complexer"
	"github.com/ph34rd/powwow/pkg/limiter"
	"github.com/ph34rd/powwow/pkg/logger"
	"github.com/ph34rd/powwow/pkg/pow"
	"github.com/ph34rd/powwow/pkg/session/manager"
	"github.com/ph34rd/powwow/pkg/stat"
	"github.com/ph34rd/powwow/pkg/wow"
)

type ServerServices struct {
	Logger    logger.Logger
	PoW       pow.PoW
	WoW       wow.WoW
	Sampler   stat.CPUSampler
	Complexer complexer.Complexer
	Manager   manager.Manager
	Limiter   limiter.StringLimiter
}

type ContextMinterFunc func(ctx context.Context, complexity uint8, prefix []byte) ([]byte, error)
