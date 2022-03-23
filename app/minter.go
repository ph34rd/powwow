package app

import (
	"context"
	"runtime"

	"golang.org/x/crypto/sha3"

	"github.com/ph34rd/powwow/pkg/pow/hashcash"
)

func mint(ctx context.Context, complexity uint8, prefix []byte) ([]byte, error) {
	iter, err := hashcash.NewFastIter(runtime.NumCPU())
	if err != nil {
		return nil, err
	}
	minter, err := hashcash.NewParallel(sha3.New256, iter, int(complexity), runtime.NumCPU())
	if err != nil {
		return nil, err
	}
	return minter.Mint(ctx, prefix)
}
