package pool

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ph34rd/powwow/pkg/server/pool/prng"
)

const poolOvergrowthKillRatio = 0.33

type workerPool struct {
	numCPU  int64
	workChs []chan func()
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(size int) (*workerPool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size out of range: %d", size)
	}
	numCPU := int64(runtime.NumCPU())
	workChs := make([]chan func(), numCPU)
	for i := range workChs {
		workChs[i] = make(chan func(), int64(size)/numCPU)
	}
	return &workerPool{
		workChs: workChs,
		numCPU:  numCPU,
	}, nil
}

func (w *workerPool) spawn(seed uint64, fn func(), workCh <-chan func(), keep bool) {
	go func() {
		if fn != nil {
			fn()
		}
		rng := prng.New(seed)
		for fn2 := range workCh {
			fn2()
			if rng.Float64() < poolOvergrowthKillRatio {
				if keep {
					w.spawn(rng.Uint64(), nil, workCh, true)
				}
				return
			}
		}
	}()
}

func (w *workerPool) Init() {
	rng := prng.New(uint64(time.Now().UnixNano()))
	for _, workCh := range w.workChs {
		for i := 0; i < cap(workCh); i++ {
			w.spawn(rng.Uint64(), nil, workCh, true)
		}
	}
}

func (w *workerPool) Go(fn func()) {
	t := time.Now().UnixNano()
	shard := t % w.numCPU
	workCh := w.workChs[shard]
	select {
	case workCh <- fn:
	default:
		w.spawn(uint64(t), fn, workCh, false)
	}
}
