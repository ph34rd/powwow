package hashcash

import (
	"context"
	"sync"
)

type Parallel struct {
	*Single
	hFn HashFunc
	cnt int

	mu      sync.Mutex
	lastSeq int
}

func NewParallel(hFn HashFunc, iter NonceIter, c, cnt int) (*Parallel, error) {
	if cnt < 1 {
		return nil, ErrHCountRange
	}
	s, err := NewSingle(hFn, iter, c)
	if err != nil {
		return nil, err
	}
	return &Parallel{Single: s, hFn: hFn, cnt: cnt, lastSeq: -1}, nil
}

func (h *Parallel) Mint(ctx context.Context, prefix []byte) ([]byte, error) {
	if h.cnt == 1 {
		return h.Single.Mint(ctx, prefix)
	}
	resCh := make(chan []byte, h.cnt-1)
	stopCh := make(chan struct{})
	for i := 0; i < h.cnt-1; i++ {
		go h.mintSeq(prefix, resCh, stopCh)
	}

	stopped := 0
	for {
		seq := h.nextSeq()
		if seq < 0 {
			// wait for other goroutines
			for {
				if stopped == h.cnt-1 {
					return nil, ErrHCollisionNotFound
				} else {
					found := <-resCh
					if found != nil {
						close(stopCh)
						return found, nil
					}
					stopped++
				}
			}
		}
	mintLoop:
		for {
			select {
			case found := <-resCh:
				if found != nil {
					close(stopCh)
					return found, nil
				}
				stopped++
			case <-ctx.Done():
				close(stopCh)
				return nil, ctx.Err()
			default:
				nonce := h.iter.Next(seq)
				if nonce == nil {
					break mintLoop
				}
				if mintOnce(h.hash, prefix, nonce, h.complexity) {
					close(stopCh)
					return nonce, nil
				}
			}
		}
	}
}

func (h *Parallel) nextSeq() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.lastSeq < h.iter.SeqSize()-1 {
		h.lastSeq++
		return h.lastSeq
	}
	return -1
}

func (h *Parallel) mintSeq(prefix []byte, resCh chan<- []byte, stopCh <-chan struct{}) {
	hs := h.hFn()
	for {
		seq := h.nextSeq()
		if seq < 0 {
			resCh <- nil
			return
		}
		for {
			nonce := h.iter.Next(seq)
			if nonce == nil {
				break
			}
			if mintOnce(hs, prefix, nonce, h.complexity) {
				resCh <- nonce
				return
			}
			select {
			case <-stopCh:
				return
			default:
			}
		}
	}
}
