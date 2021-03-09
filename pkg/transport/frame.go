package transport

import "encoding/binary"

const frameHeaderSize = 2
const frameHeaderPldSizeSize = 4
const frameHeaderPldTypeSize = 2

// FrameHeader represents single frame in the stream.
// |1   |2   |3   |4-8   |8-16 |16-48  |48-64  |
// |rsv1|rsv2|rsv3|opCode|flags|pldSize|pldType|
type FrameHeader [frameHeaderSize + frameHeaderPldSizeSize + frameHeaderPldTypeSize]byte

type Flags uint8

const (
	FlagHasPayload     = 0b10000000
	FlagHasPayloadType = 0b01000000
)

func (f Flags) IsSet(cf Flags) bool {
	return f&cf != 0
}

func (f *Flags) Set(cf Flags) {
	*f |= cf
}

// OpCode represents operation code.
type OpCode uint8

const (
	OpData  OpCode = 0x1
	OpClose OpCode = 0x8
	OpPing  OpCode = 0x9
	OpPong  OpCode = 0xa
)

// IsControl checks whether the c is control operation code.
func (c OpCode) IsControl() bool {
	return c&0x8 != 0
}

// IsData checks whether the c is data operation code.
func (c OpCode) IsData() bool {
	return c&0x8 == 0
}

// IsValid checks whether the c is defined operation code.
func (c OpCode) IsValid() bool {
	return c == OpData || (c >= OpClose && c <= OpPong)
}

// CheckRsv checks reserved fields.
func (f FrameHeader) CheckRsv() bool {
	return f[0]&0xf0 == 0
}

// OpCode is the getter of OpCode.
func (f FrameHeader) OpCode() OpCode {
	return OpCode(f[0] & 0x0f)
}

// PayloadSize is the getter of PayloadSize.
func (f FrameHeader) PayloadSize() uint32 {
	return binary.BigEndian.Uint32(f[2:])
}

// Flags is the getter of Flags.
func (f FrameHeader) Flags() Flags {
	return Flags(f[1])
}

// PayloadType is the getter of PayloadType.
func (f FrameHeader) PayloadType() uint16 {
	return binary.BigEndian.Uint16(f[6:])
}

// SetFlags is the setter of Flags.
func (f *FrameHeader) SetFlags(cf Flags) {
	f[1] |= uint8(cf)
}

// SetOpCode is the setter of OpCode.
func (f *FrameHeader) SetOpCode(oc OpCode) {
	f[0] = (f[0] & 0xf0) | (uint8(oc) & 0x0f)
}

// SetPayloadSize is the setter of PayloadSize.
func (f *FrameHeader) SetPayloadSize(sz uint32) {
	binary.BigEndian.PutUint32(f[2:], sz)
}

// SetPayloadType is the setter of PayloadType.
func (f *FrameHeader) SetPayloadType(typ uint16) {
	binary.BigEndian.PutUint16(f[6:], typ)
}

// MakeDataFrameHeader forms single frame with OpData operation code.
func MakeDataFrameHeader(sz uint32, typ uint16) FrameHeader {
	var frame FrameHeader
	frame.SetOpCode(OpData)
	frame.SetPayloadSize(sz)
	frame.SetFlags(FlagHasPayload)
	if typ != 0 {
		frame.SetPayloadType(typ)
		frame.SetFlags(FlagHasPayloadType)
	}
	return frame
}

var (
	pingFrame  = makeControlFrame(OpPing)
	pongFrame  = makeControlFrame(OpPong)
	closeFrame = makeControlFrame(OpClose)
)

func makeControlFrame(code OpCode) FrameHeader {
	var frame FrameHeader
	frame.SetOpCode(code)
	return frame
}
