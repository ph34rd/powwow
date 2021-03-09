package hashcash

import (
	"context"
	"hash"
)

type Single struct {
	hash       hash.Hash
	iter       NonceIter
	complexity int
}

func NewSingle(hFn HashFunc, iter NonceIter, c int) (*Single, error) {
	h := hFn()
	if c >= h.Size()*8 {
		return nil, ErrHComplexity
	}
	return &Single{hash: h, iter: iter, complexity: c}, nil
}

func (h *Single) Mint(ctx context.Context, prefix []byte) ([]byte, error) {
	var nonce []byte
	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		nonce = h.iter.Next(i)
		if nonce == nil {
			break
		}
		if mintOnce(h.hash, prefix, nonce, h.complexity) {
			return nonce, nil
		}
		i++
		if i > h.iter.SeqSize()-1 {
			i = 0
		}
	}
	return nil, ErrHCollisionNotFound
}
