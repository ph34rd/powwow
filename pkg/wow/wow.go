package wow

import (
	"context"
	"sync/atomic"
)

type WoW interface {
	GetNext(ctx context.Context) (string, error)
}

type inMemory struct {
	data []string
	cnt  int32
}

func NewInMemory() *inMemory {
	return &inMemory{data: data, cnt: -1}
}

func (m *inMemory) nextIdx() int {
	v := int(atomic.AddInt32(&m.cnt, 1))
	if v < 0 {
		v = v * -1
	}
	return v % len(m.data)
}

func (m *inMemory) GetNext(_ context.Context) (string, error) {
	return m.data[m.nextIdx()], nil
}
