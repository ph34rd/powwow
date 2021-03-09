// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: session.proto

package pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ServerHandshake struct {
	Challenge  []byte `protobuf:"bytes,1,opt,name=challenge,proto3" json:"challenge,omitempty"`
	Complexity uint32 `protobuf:"varint,2,opt,name=complexity,proto3" json:"complexity,omitempty"`
}

func (m *ServerHandshake) Reset()         { *m = ServerHandshake{} }
func (m *ServerHandshake) String() string { return proto.CompactTextString(m) }
func (*ServerHandshake) ProtoMessage()    {}
func (*ServerHandshake) Descriptor() ([]byte, []int) {
	return fileDescriptor_3a6be1b361fa6f14, []int{0}
}
func (m *ServerHandshake) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ServerHandshake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ServerHandshake.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ServerHandshake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerHandshake.Merge(m, src)
}
func (m *ServerHandshake) XXX_Size() int {
	return m.Size()
}
func (m *ServerHandshake) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerHandshake.DiscardUnknown(m)
}

var xxx_messageInfo_ServerHandshake proto.InternalMessageInfo

func (m *ServerHandshake) GetChallenge() []byte {
	if m != nil {
		return m.Challenge
	}
	return nil
}

func (m *ServerHandshake) GetComplexity() uint32 {
	if m != nil {
		return m.Complexity
	}
	return 0
}

type ClientHandshake struct {
	Nonce []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *ClientHandshake) Reset()         { *m = ClientHandshake{} }
func (m *ClientHandshake) String() string { return proto.CompactTextString(m) }
func (*ClientHandshake) ProtoMessage()    {}
func (*ClientHandshake) Descriptor() ([]byte, []int) {
	return fileDescriptor_3a6be1b361fa6f14, []int{1}
}
func (m *ClientHandshake) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClientHandshake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClientHandshake.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClientHandshake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientHandshake.Merge(m, src)
}
func (m *ClientHandshake) XXX_Size() int {
	return m.Size()
}
func (m *ClientHandshake) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientHandshake.DiscardUnknown(m)
}

var xxx_messageInfo_ClientHandshake proto.InternalMessageInfo

func (m *ClientHandshake) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

type WoWRequest struct {
}

func (m *WoWRequest) Reset()         { *m = WoWRequest{} }
func (m *WoWRequest) String() string { return proto.CompactTextString(m) }
func (*WoWRequest) ProtoMessage()    {}
func (*WoWRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3a6be1b361fa6f14, []int{2}
}
func (m *WoWRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WoWRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WoWRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WoWRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WoWRequest.Merge(m, src)
}
func (m *WoWRequest) XXX_Size() int {
	return m.Size()
}
func (m *WoWRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WoWRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WoWRequest proto.InternalMessageInfo

type WoWResponse struct {
	Wow string `protobuf:"bytes,1,opt,name=wow,proto3" json:"wow,omitempty"`
}

func (m *WoWResponse) Reset()         { *m = WoWResponse{} }
func (m *WoWResponse) String() string { return proto.CompactTextString(m) }
func (*WoWResponse) ProtoMessage()    {}
func (*WoWResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3a6be1b361fa6f14, []int{3}
}
func (m *WoWResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WoWResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WoWResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WoWResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WoWResponse.Merge(m, src)
}
func (m *WoWResponse) XXX_Size() int {
	return m.Size()
}
func (m *WoWResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WoWResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WoWResponse proto.InternalMessageInfo

func (m *WoWResponse) GetWow() string {
	if m != nil {
		return m.Wow
	}
	return ""
}

func init() {
	proto.RegisterType((*ServerHandshake)(nil), "powwow.session.ServerHandshake")
	proto.RegisterType((*ClientHandshake)(nil), "powwow.session.ClientHandshake")
	proto.RegisterType((*WoWRequest)(nil), "powwow.session.WoWRequest")
	proto.RegisterType((*WoWResponse)(nil), "powwow.session.WoWResponse")
}

func init() { proto.RegisterFile("session.proto", fileDescriptor_3a6be1b361fa6f14) }

var fileDescriptor_3a6be1b361fa6f14 = []byte{
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2e,
	0xce, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b, 0xc8, 0x2f, 0x2f, 0xcf,
	0x2f, 0xd7, 0x83, 0x8a, 0x4a, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0xa5, 0xf4, 0x41, 0x2c, 0x88,
	0x2a, 0x25, 0x7f, 0x2e, 0xfe, 0xe0, 0xd4, 0xa2, 0xb2, 0xd4, 0x22, 0x8f, 0xc4, 0xbc, 0x94, 0xe2,
	0x8c, 0xc4, 0xec, 0x54, 0x21, 0x19, 0x2e, 0xce, 0xe4, 0x8c, 0xc4, 0x9c, 0x9c, 0xd4, 0xbc, 0xf4,
	0x54, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x84, 0x80, 0x90, 0x1c, 0x17, 0x57, 0x72, 0x7e,
	0x6e, 0x41, 0x4e, 0x6a, 0x45, 0x66, 0x49, 0xa5, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x6f, 0x10, 0x92,
	0x88, 0x92, 0x3a, 0x17, 0xbf, 0x73, 0x4e, 0x66, 0x6a, 0x5e, 0x09, 0xc2, 0x40, 0x11, 0x2e, 0xd6,
	0xbc, 0xfc, 0xbc, 0x64, 0x98, 0x61, 0x10, 0x8e, 0x12, 0x0f, 0x17, 0x57, 0x78, 0x7e, 0x78, 0x50,
	0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x92, 0x3c, 0x17, 0x37, 0x98, 0x57, 0x5c, 0x90, 0x9f, 0x57,
	0x9c, 0x2a, 0x24, 0xc0, 0xc5, 0x5c, 0x9e, 0x5f, 0x0e, 0xd6, 0xc0, 0x19, 0x04, 0x62, 0x3a, 0xc9,
	0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb,
	0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x14, 0x53, 0x41, 0x52, 0x12, 0x1b,
	0xd8, 0x37, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x9c, 0x52, 0xf2, 0x04, 0x01, 0x00,
	0x00,
}

func (m *ServerHandshake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ServerHandshake) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ServerHandshake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Complexity != 0 {
		i = encodeVarintSession(dAtA, i, uint64(m.Complexity))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Challenge) > 0 {
		i -= len(m.Challenge)
		copy(dAtA[i:], m.Challenge)
		i = encodeVarintSession(dAtA, i, uint64(len(m.Challenge)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ClientHandshake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientHandshake) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClientHandshake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Nonce) > 0 {
		i -= len(m.Nonce)
		copy(dAtA[i:], m.Nonce)
		i = encodeVarintSession(dAtA, i, uint64(len(m.Nonce)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WoWRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WoWRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WoWRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *WoWResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WoWResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WoWResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Wow) > 0 {
		i -= len(m.Wow)
		copy(dAtA[i:], m.Wow)
		i = encodeVarintSession(dAtA, i, uint64(len(m.Wow)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSession(dAtA []byte, offset int, v uint64) int {
	offset -= sovSession(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ServerHandshake) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Challenge)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	if m.Complexity != 0 {
		n += 1 + sovSession(uint64(m.Complexity))
	}
	return n
}

func (m *ClientHandshake) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Nonce)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	return n
}

func (m *WoWRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *WoWResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Wow)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	return n
}

func sovSession(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSession(x uint64) (n int) {
	return sovSession(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ServerHandshake) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSession
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ServerHandshake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ServerHandshake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Challenge", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Challenge = append(m.Challenge[:0], dAtA[iNdEx:postIndex]...)
			if m.Challenge == nil {
				m.Challenge = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Complexity", wireType)
			}
			m.Complexity = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Complexity |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSession(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSession
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ClientHandshake) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSession
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClientHandshake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientHandshake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nonce = append(m.Nonce[:0], dAtA[iNdEx:postIndex]...)
			if m.Nonce == nil {
				m.Nonce = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSession(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSession
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WoWRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSession
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WoWRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WoWRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipSession(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSession
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WoWResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSession
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WoWResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WoWResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Wow", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Wow = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSession(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSession
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSession(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSession
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSession
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSession
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSession
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSession
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSession
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSession        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSession          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSession = fmt.Errorf("proto: unexpected end of group")
)