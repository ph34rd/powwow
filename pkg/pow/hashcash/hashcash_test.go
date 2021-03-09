package hashcash

import (
	"context"
	"testing"

	"golang.org/x/crypto/sha3"
)

func Test_Hashcash(t *testing.T) {
	challenge := []byte("0123456789")
	complexity := 20

	iter, err := NewFastIter(1)
	if err != nil {
		t.Fatalf("expected nil error: got: %v", err)
	}
	h, err := NewParallel(sha3.New256, iter, complexity, 1)
	if err != nil {
		t.Fatalf("expected nil error: got: %v", err)
	}

	nonce, err := h.Mint(context.Background(), challenge)
	if err != nil {
		t.Fatalf("expected nil error: got: %v", err)
	}

	err = NewValidator(sha3.New256).Validate(challenge, nonce, complexity)
	if err != nil {
		t.Fatalf("expected nil error: got: %v", err)
	}
}
