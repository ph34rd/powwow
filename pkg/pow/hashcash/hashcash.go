package hashcash

import (
	"context"
	"errors"
	"hash"
)

var (
	ErrHComplexity        = errors.New("hashcash: complexity is too high")
	ErrHCollisionNotFound = errors.New("hashcash: collision not found")
	ErrHNotValidated      = errors.New("hashcash: nonce not validated")
	ErrHCountRange        = errors.New("hashcash: count is out of range")
)

type HashFunc func() hash.Hash

type NonceIter interface {
	Next(seq int) []byte
	SeqSize() int
}

type Minter interface {
	Mint(ctx context.Context, prefix []byte) ([]byte, error)
}

type Validator interface {
	Validate(prefix, nonce []byte, complexity int) error
}
