package hashcash

import (
	"encoding/binary"
)

const fastIterNonceSize = 8

type fastIterState struct {
	start   uint64
	current uint64
	stop    uint64
	value   [fastIterNonceSize]byte
}

type FastIter struct {
	seqS []fastIterState
}

func NewFastIter(seq int) (*FastIter, error) {
	if seq < 1 {
		return nil, ErrHCountRange
	}
	if seq == 1 {
		return &FastIter{seqS: []fastIterState{{stop: 0xffffffffffffffff}}}, nil
	}
	seqSize := 0xffffffffffffffff / uint64(seq)
	var start, stop uint64
	seqS := make([]fastIterState, seq)
	for i := 0; i < seq; i++ {
		seqS[i].start = start
		seqS[i].current = start
		stop = start + seqSize - 1
		seqS[i].stop = stop
		start += seqSize
	}
	seqS[seq-1].stop = 0xffffffffffffffff
	return &FastIter{seqS: seqS}, nil
}

func (it *FastIter) Next(seq int) []byte {
	cur := it.seqS[seq].current
	if cur > it.seqS[seq].stop || cur < it.seqS[seq].start {
		return nil
	}
	binary.LittleEndian.PutUint64(it.seqS[seq].value[:], cur)
	it.seqS[seq].current++
	return it.seqS[seq].value[:]
}

func (it *FastIter) SeqSize() int {
	return len(it.seqS)
}
