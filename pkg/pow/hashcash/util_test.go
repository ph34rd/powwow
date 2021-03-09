package hashcash

import (
	"strconv"
	"testing"
)

func Test_countLeadingZeroBits(t *testing.T) {
	tests := []struct {
		in     []byte
		expect int
	}{
		{in: []byte{0, 0b00000001}, expect: 15},
		{in: []byte{0, 0b00000011}, expect: 14},
		{in: []byte{0, 0b00000111}, expect: 13},
		{in: []byte{0, 0b00001111}, expect: 12},
		{in: []byte{0, 0b00011111}, expect: 11},
		{in: []byte{0, 0b00111111}, expect: 10},
		{in: []byte{0, 0b01111111}, expect: 9},
		{in: []byte{0, 0xff}, expect: 8},
		{in: []byte{0}, expect: 8},
		{in: []byte{0xff}, expect: 0},
		{in: []byte{0xf2, 0x0c, 0x0f}, expect: 0},
		{in: nil, expect: 0},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cnt := countLeadingZeroBits(test.in)
			if cnt != test.expect {
				t.Fatalf("expected result: %d, got: %d", test.expect, cnt)
			}
		})
	}
}
