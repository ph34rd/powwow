package hashcash

import (
	"hash"
)

// countLeadingZeroBits counts leading zero bits in b.
func countLeadingZeroBits(b []byte) (cnt int) {
	for _, x := range b {
		switch {
		case x == 0:
			cnt += 8
		case x&128 != 0:
			return
		case x&64 != 0:
			cnt++
			return
		case x&32 != 0:
			cnt += 2
			return
		case x&16 != 0:
			cnt += 3
			return
		case x&8 != 0:
			cnt += 4
			return
		case x&4 != 0:
			cnt += 5
			return
		case x&2 != 0:
			cnt += 6
			return
		case x&1 != 0:
			cnt += 7
			return
		}
	}
	return
}

// mintOnce checks the solution.
func mintOnce(hs hash.Hash, prefix, nonce []byte, complexity int) bool {
	hs.Reset()
	hs.Write(prefix)
	hs.Write(nonce[:])
	return countLeadingZeroBits(hs.Sum(nil)) >= complexity
}
