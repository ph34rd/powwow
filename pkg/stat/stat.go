package stat

import (
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

// CPUUsage contains process CPU times.
type CPUUsage struct {
	System time.Duration
	User   time.Duration
}

type CPUSampler interface {
	GetCPUUsageFraction() float64
	Shutdown()
}

type SamplerImpl struct {
	logger *zap.Logger
	numCPU int
	period time.Duration

	once   sync.Once
	stopCh chan struct{}

	mu       sync.RWMutex
	fraction float64
}

func NewCPUSampler(lg *zap.Logger, period time.Duration) *SamplerImpl {
	ret := &SamplerImpl{
		logger: lg,
		numCPU: runtime.NumCPU(),
		period: period,
		stopCh: make(chan struct{}),
	}
	go ret.runSampler()
	return ret
}

func (s *SamplerImpl) runSampler() {
	prevTime := time.Now()
	prev, err := GetCPUUsage()
	if err != nil {
		s.logger.Error("get cpu sample error", zap.Error(err))
	}
	t := time.NewTicker(s.period)
	for {
		select {
		case curTime := <-t.C:
			cur, err := GetCPUUsage()
			if err != nil {
				s.logger.Error("get cpu sample error", zap.Error(err))
			}
			fraction := calcCPuUsagePercent(prev, cur, curTime.Sub(prevTime), s.numCPU)
			s.mu.Lock()
			s.fraction = fraction
			s.mu.Unlock()
			prev = cur
		case <-s.stopCh:
			t.Stop()
			return
		}
	}
}

func calcCPuUsagePercent(prev, cur CPUUsage, elapsed time.Duration, nCPU int) float64 {
	total := elapsed.Seconds() * float64(nCPU)
	var userDiff, systemDiff time.Duration
	if prev.User != 0 && cur.User > prev.User {
		userDiff = cur.User - prev.User
	}
	if prev.System != 0 && cur.System > prev.System {
		systemDiff = cur.System - prev.System
	}
	if total < 0 {
		return 0
	}
	return (userDiff.Seconds() + systemDiff.Seconds()) / total
}

func (s *SamplerImpl) GetCPUUsageFraction() (fraction float64) {
	s.mu.RLock()
	fraction = s.fraction
	s.mu.RUnlock()
	return
}

func (s *SamplerImpl) Shutdown() {
	s.once.Do(func() { close(s.stopCh) })
}
