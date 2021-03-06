package prng

import (
	"math"
	"math/bits"
	"unsafe"
)

// set true for some generator tests
// const twistedUint64 = true
const twistedUint64 = false

// A Xoro with a xoroshiro prng implements a 64-bit generator with 128-bit state.
// A Xoro is the den of the xoroshiros holding their two 64-bit eggs.
type Xoro struct {
	s0, s1 uint64
}

// NewXoro returns a new xoroshiro128 generator seeded by the seed.
// Float64 uses xoroshiro128+ and Uint64 xoroshiro128**. Both xoroshiros update
// a Xoro in the same way (same linear engine) and we can use a single Xoro for both
// functions without interfering random stream properties.
func NewXoro(seed uint64) Xoro {
	x := Xoro{}
	x.Seed(seed)
	return x
}

// Seed seeds a xorohiro128 generator by seed using splitMix64. Any seed is ok.
func (x *Xoro) Seed(seed uint64) {
	x.s0 = Splitmix(&seed)
	x.s1 = Splitmix(&seed)
}

// NextXoro returns the next xoroshiro128 from Outlet. Each generator has
// 2^64 long random streams, which is not overlapping with other generators streams.
// NextXoro is safe for concurrent use by multiple goroutines.
func (s *Outlet) NextXoro() Xoro {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.xoro.Jump()
	return s.xoro
}

// NewXoroSlice returns a slice of n xoroshiro128 generators with non-overlapping 2^64
// long random streams. First generator is seeded by seed.
func NewXoroSlice(n int, seed uint64) []Xoro {
	s := make([]Xoro, n)
	s[0].Seed(seed)
	for i := 1; i < n; i++ {
		s[i] = s[i-1]
		s[i].Jump()
	}
	return s
}

// Uint64 returns a  pseudo-random uint64. Uint64 is xoroshiro128**.
func (x *Xoro) Uint64() (next uint64) {

	next = bits.RotateLeft64(x.s0*5, 7) * 9
	*x = x.NextState()
	if twistedUint64 { // for testing
		next = twisted(next)
	}
	return
}

func twisted(next uint64) uint64 {
	shift := (next >> 10) & 63
	if next%2 == 0 {
		return 0
	}
	if next%3 == 0 {
		return next >> shift
	}
	if next%5 == 0 {
		return next << shift
	}
	if next%7 == 0 {
		// this catches rounding differencies
		return ((1 << 64) - 1) >> shift
	}
	if next%11 == 0 {
		return ((1 << 64) - 1) << shift
	}
	return next
}

// Xoroshiro128plus is xoroshiro128+
func (x *Xoro) Xoroshiro128plus() (next uint64) {

	next = x.s0 + x.s1
	*x = x.NextState()
	return
}

// NextState returns the next Xoro state of the xoroshiro128+/** linear engine.
func (x Xoro) NextState() Xoro {
	// gc compiler detects similar expressions if given in parentheses

	return Xoro{
		s0: bits.RotateLeft64(x.s0, 24) ^ (x.s0 ^ x.s1) ^ ((x.s0 ^ x.s1) << 16),
		s1: bits.RotateLeft64(x.s0^x.s1, 37),
	}
}

// WriteState writes the current state of the generator x to b.
// WriteState without allocation is faster than State().
func (x *Xoro) WriteState(b []byte) {
	if len(b) < XoroStateSize {
		panic("WriteState: byte slice too short")
	}
	// This expects a little endian cpu, eg. all amd64.
	*(*uint64)(unsafe.Pointer(&b[0])) = bits.ReverseBytes64(x.s0)
	*(*uint64)(unsafe.Pointer(&b[8])) = bits.ReverseBytes64(x.s1)
}

// State returns the current state of the generator x as []byte.
func (x *Xoro) State() []byte {
	var b [XoroStateSize]byte

	*(*uint64)(unsafe.Pointer(&b[0])) = bits.ReverseBytes64(x.s0)
	*(*uint64)(unsafe.Pointer(&b[8])) = bits.ReverseBytes64(x.s1)
	return b[:]
}

// ReadState reads the state of the generator x from b []byte.
func (x *Xoro) ReadState(b []byte) {
	if len(b) < XoroStateSize {
		panic("ReadState: byte slice too short")
	}
	x.s0 = bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(&b[0])))
	x.s1 = bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(&b[8])))
}

// Float64 returns a uniformly distributed pseudo-random float64 from [0, 1).
// The distribution is  2^53 evenly spaced floats with spacing 2^-53.
func (x *Xoro) Float64() float64 {

	return float64(x.Xoroshiro128plus()>>11) / (1 << 53)
}

// Float64_64 returns a uniformly distributed pseudo-random float64 from [0, 1).
// The distribution includes all floats in [2^-12, 1) and 2^52 evenly spaced
// floats in [0, 2^-12) with spacing 2^-64.
func (x *Xoro) Float64_64() float64 {

	u := x.Uint64()
	if u == 0 {
		return 0
	} // without this the smallest returned is 2^-65
	z := uint64(bits.LeadingZeros64(u)) + 1
	return math.Float64frombits((1023-z)<<52 | u<<z>>12)
}

// Float64_64R returns a uniformly distributed pseudo-random float64 from [0, 1]
// using rounding. The distribution includes all rounded floats in [2^-11, 1]
// and 2^53 evenly spaced floats in [0, 2^-11) with spacing 2^-64.
func (x *Xoro) Float64_64R() float64 {

	u := x.Uint64()
	if u == 0 {
		return 0
	}
	z := uint64(bits.LeadingZeros64(u)) + 1
	return math.Float64frombits((((1023-z)<<53 | u<<z>>11) + 1) >> 1)
}

// Float64full returns a uniformly distributed pseudo-random float64 from [0, 1).
// The distribution includes all floats in [0, 1).
// Float64full is equivalent to Float64Bisect in truncate mode.
func (x *Xoro) Float64full() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u)) + 1
	if z <= 12 { // 99.975% of cases
		return math.Float64frombits((1023-z)<<52 | u<<z>>12)
	}
	z--
	exp := uint64(0)
	for u == 0 {
		u = x.Uint64()
		z = uint64(bits.LeadingZeros64(u))
		exp += 64
		if exp+z >= 1074 {
			return 0
		}
	}
	u = u<<z | x.Uint64()>>(64-z)
	exp += z
	if exp < 1022 {
		return math.Float64frombits((1022-exp)<<52 | u<<1>>12)
	}
	return math.Float64frombits(u >> (exp - 1022) >> 12) // 2^52 subnormal floats
}

// Float64_128 --
func (x *Xoro) Float64_128() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u)) + 1
	if z <= 12 { // 99.975% of cases
		return math.Float64frombits((1023-z)<<52 | u<<z>>12)
	}
	z--
	exp := z
	if u == 0 {
		u = x.Uint64()
		z = uint64(bits.LeadingZeros64(u))
		exp = 64 + z
		if exp >= 128 {
			return 0
		}
	}
	u = u<<z | x.Uint64()>>(64-z)
	return math.Float64frombits((1022-exp)<<52 | u<<1>>12)
}

// Float64fullR returns a uniformly distributed pseudo-random float64 from [0, 1]
// using rounding. The distribution includes all floats in [0, 1].
// Float64fullR is equivalent to Float64Bisect in rounding mode.
func (x *Xoro) Float64fullR() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u)) + 1
	if z <= 11 { // 99.95% of cases
		return math.Float64frombits((((1023-z)<<53 | u<<z>>11) + 1) >> 1)
	}
	z--
	exp := uint64(0)
	for u == 0 {
		u = x.Uint64()
		z = uint64(bits.LeadingZeros64(u))
		exp += 64
		if exp+z > 1074 {
			return 0
		}
	}
	u = u<<z | x.Uint64()>>(64-z)
	exp += z
	if exp < 1022 {
		return math.Float64frombits((((1022-exp)<<53 | u<<1>>11) + 1) >> 1)
	}
	return math.Float64frombits((u>>(exp-1022)>>11 + 1) >> 1)
}

// twoToMinus(n) returns 2^-n as a float64.
func twoToMinus(n uint64) float64 {
	n = (1023 - n) << 52
	return *(*float64)(unsafe.Pointer(&n))
}

// ldexp(f, exp) returns f * 2^-exp as a float64.
func ldexp(f float64, exp uint64) float64 {
	if exp > 1022 {
		f *= 0x1p-1022
		exp -= 1022
	}
	return f * twoToMinus(exp)
}

// Float64_117 returns a uniformly distributed pseudo-random float64 from [0, 1).
// The distribution includes all floats in [2^-65, 1) and 2^52 evenly spaced
// floats in [0, 2^-65) with spacing 2^-117.
func (x *Xoro) Float64_117() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u)) + 1
	if z <= 12 {
		return math.Float64frombits((1023-z)<<52 | u<<z>>12)
	}
	z--
	u = u<<z | x.Uint64()>>(64-z)
	return float64(u>>11) * twoToMinus(53+z)
}

// Float64_117R returns a uniformly distributed pseudo-random float64 from
// [0, 1] using rounding. The distribution includes all floats in [2^-65, 1]
// and 2^52 evenly spaced floats in [0, 2^-65) with spacing 2^-117.
func (x *Xoro) Float64_117R() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u)) + 1
	if z <= 11 {
		return math.Float64frombits((((1023-z)<<53 | u<<z>>11) + 1) >> 1)
	}
	z--
	u = u<<z | x.Uint64()>>(64-z)
	return float64((u>>10+1)>>1) * twoToMinus(53+z)
}

// RandomReal returns a uniformly distributed pseudo-random float64 from [0, 1].
// The distribution includes all floats, but may miss very few
// subnormal floats in in [0, 2^-1022).
// http://prng.di.unimi.it/random_real.c
// RandomReal is equivalent to Float64Bisect in rounding mode in [2^-1024, 1].
func (x *Xoro) RandomReal() float64 {

	u := x.Uint64()
	z := uint64(bits.LeadingZeros64(u))
	exp := uint64(64)
	for u == 0 {
		u = x.Uint64()
		z = uint64(bits.LeadingZeros64(u))
		exp += 64
		if exp+z > 1074+64 {
			return 0
		}

	}
	u = u<<z | x.Uint64()>>(64-z)
	return ldexp(float64(u|1), exp+z)
}

// Float64Bisect returns a uniformly distributed pseudo-random float64 value in [0, 1).
// If round is true, rounding is applied and the range is [0, 1].
// All floats, normal and subnormal, are included.
func (x *Xoro) Float64Bisect(round bool) float64 {

	left, mean, right := 0.0, 0.5, 1.0
	for {
		u := x.Uint64()
		for b := 0; b < 64; b++ {

			if u&(1<<63) != 0 { // evaluate the leftmost bit of u
				left = mean // '1' bit -> take the right half, big numbers
			} else {
				right = mean // '0' bit -> take the left half, small numbers
			}
			u <<= 1
			mean = (left + right) / 2
			if mean == left || mean == right { // check if left and right are adjacent floats
				if !round {
					return left // no rounding
				}
				if b == 63 { // must have one rounding bit
					u = x.Uint64()
				}
				if u&(1<<63) != 0 { // evaluate the rounding bit
					return right // '1' bit -> round up
				}
				return left // '0' bit -> round down
			}
		}
	}
}
