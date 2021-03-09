package transport

import (
	"testing"
)

func Test_frameHeader_SetOpCode(t *testing.T) {
	var f FrameHeader
	f[0] = 0b01010101
	f.SetOpCode(OpCode(0xf))
	if f[0] != 0b01011111 {
		t.Fatalf("expected result: %b, got: %b", 0b01011111, f[0])
	}
	if f.OpCode() != 0x0f {
		t.Fatalf("expected result: %b, got: %b", 0x0f, f.OpCode())
	}
}

func Test_frameHeader_SetPayloadSize(t *testing.T) {
	var f FrameHeader
	f[0] = 0xff
	f.SetPayloadSize(0xffff)
	if f[2] != 0x00 || f[3] != 0x00 || f[4] != 0xff || f[5] != 0xff {
		t.Fatalf("expected result: %x, got: %x", 0x0000ffff, f[1:])
	}
	if f.PayloadSize() != 0xffff {
		t.Fatalf("expected result: %x, got: %x", 0xffff, f.PayloadSize())
	}
}
