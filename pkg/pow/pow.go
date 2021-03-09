package pow

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/ph34rd/powwow/pkg/pow/hashcash"
)

const (
	challengeSize = 16
	dateSize      = 8
)

type PoW interface {
	Challenge() ([]byte, error)
	Verify(challenge, nonce []byte, complexity uint8) error
}

type Svc struct {
	hFn hashcash.HashFunc
}

func New(hFn hashcash.HashFunc) *Svc {
	return &Svc{hFn: hFn}
}

func (s Svc) Challenge() ([]byte, error) {
	buf := make([]byte, challengeSize+dateSize)
	_, err := rand.Read(buf[:challengeSize])
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(buf[challengeSize:], uint64(time.Now().UnixNano()))
	return buf, nil
}

func (s *Svc) Verify(challenge, nonce []byte, c uint8) error {
	return hashcash.NewValidator(s.hFn).Validate(challenge, nonce, int(c))
}
