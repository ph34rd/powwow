package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultBurst = 10
	defaultEvery = 10 * time.Second
)

type StringLimiter interface {
	Add(string) *rate.Limiter
	Get(string) *rate.Limiter
}

type Impl struct {
	ips   map[string]*rate.Limiter
	mu    sync.RWMutex
	burst int
	rate  rate.Limit
}

func NewStringLimiter() *Impl {
	return &Impl{
		ips:   make(map[string]*rate.Limiter),
		burst: defaultBurst,
		rate:  rate.Every(defaultEvery),
	}
}

func (i *Impl) Add(s string) *rate.Limiter {
	limiter := rate.NewLimiter(i.rate, i.burst)
	i.mu.Lock()
	i.ips[s] = limiter
	i.mu.Unlock()
	return limiter
}

func (i *Impl) Get(s string) *rate.Limiter {
	i.mu.RLock()
	limiter, ok := i.ips[s]
	if !ok {
		i.mu.RUnlock()
		return i.Add(s)
	}
	i.mu.RUnlock()
	return limiter
}
