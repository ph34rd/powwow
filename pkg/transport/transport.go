package transport

import (
	"errors"
	"io"
)

const defaultMaxFrameSize = 65536

var (
	ErrFrameTooLarge     = errors.New("tlv: frame too large")
	ErrMalformedHeader   = errors.New("tlv: malformed header")
	ErrUnexpectedEOF     = errors.New("tlv: unexpected EOF")
	ErrNotControlFrame   = errors.New("tlv: not a control frame")
	ErrThresholdExceeded = errors.New("tlv: threshold exceeded")
	ErrUncompletedWrite  = errors.New("tlv: uncompleted write")
	ErrClose             = errors.New("tlv: close")
)

type (
	DataFrameReaderFunc func(size uint32, typeID uint16, r io.Reader) error
	DataFrameWriterFunc func(w io.Writer) error
)

type FrameReader interface {
	NextReader(fn DataFrameReaderFunc) (err error)
}

type FrameWriter interface {
	NextWriter(size uint32, typeID uint16, fn DataFrameWriterFunc) error
}

type FrameCloser interface {
	Close() error
}

type FramePinger interface {
	Ping() error
}

type Transport interface {
	FrameReader
	FrameWriter
	FrameCloser
	FramePinger
}

type TLVTransport struct {
	r            io.Reader
	w            io.Writer
	maxFrameSize uint32
}

func NewTLVTransport(r io.Reader, w io.Writer, maxFrameSize int) *TLVTransport {
	if maxFrameSize <= 0 {
		maxFrameSize = defaultMaxFrameSize
	}
	return &TLVTransport{
		r:            r,
		w:            w,
		maxFrameSize: uint32(maxFrameSize),
	}
}

func ReadFrameHeader(r io.Reader) (FrameHeader, error) {
	var frame FrameHeader
	_, err := io.ReadFull(r, frame[:frameHeaderSize])
	if err != nil {
		return FrameHeader{}, err
	}
	if !frame.CheckRsv() {
		return FrameHeader{}, ErrMalformedHeader
	}
	if frame.OpCode().IsControl() {
		if frame.Flags().IsSet(FlagHasPayload) || frame.Flags().IsSet(FlagHasPayloadType) {
			return FrameHeader{}, ErrMalformedHeader
		}
		return frame, nil
	}
	if !frame.Flags().IsSet(FlagHasPayload) {
		return FrameHeader{}, ErrMalformedHeader
	}
	toRead := frameHeaderPldSizeSize
	if frame.Flags().IsSet(FlagHasPayloadType) {
		toRead += frameHeaderPldTypeSize
	} else {
		frame.SetPayloadType(0)
	}
	_, err = io.ReadFull(r, frame[frameHeaderSize:frameHeaderSize+toRead])
	if err != nil {
		return FrameHeader{}, err
	}
	return frame, nil
}

func WriteFrameHeader(frame FrameHeader, w io.Writer) error {
	if !frame.CheckRsv() {
		return ErrMalformedHeader
	}
	toWrite := frameHeaderSize
	if frame.OpCode().IsControl() {
		if frame.Flags().IsSet(FlagHasPayload) || frame.Flags().IsSet(FlagHasPayloadType) {
			return ErrMalformedHeader
		}
	} else {
		if !frame.Flags().IsSet(FlagHasPayload) {
			return ErrMalformedHeader
		}
		toWrite += frameHeaderPldSizeSize
		if frame.Flags().IsSet(FlagHasPayloadType) {
			toWrite += frameHeaderPldTypeSize
		}
	}
	_, err := w.Write(frame[:toWrite])
	return err
}

func (t *TLVTransport) NextReader(fn DataFrameReaderFunc) (retErr error) {
	rHeader, err := ReadFrameHeader(t.r)
	if err != nil {
		return err
	}
	if rHeader.OpCode().IsControl() {
		return t.handleControlFrame(rHeader.OpCode())
	}
	if rHeader.PayloadSize() > t.maxFrameSize {
		return ErrFrameTooLarge
	}
	payloadReader := io.Reader(nullReader{})
	if rHeader.PayloadSize() != 0 {
		lr := &io.LimitedReader{R: t.r, N: int64(rHeader.PayloadSize())}
		defer func() {
			err := discardReader(lr)
			if err != nil {
				retErr = err
			}
		}()
		payloadReader = lr
	}
	return fn(rHeader.PayloadSize(), rHeader.PayloadType(), payloadReader)
}

func (t *TLVTransport) handleControlFrame(code OpCode) error {
	switch code {
	case OpPing:
		err := WriteFrameHeader(pongFrame, t.w)
		if err != nil {
			return err
		}
		return flushWriter(t.w)
	case OpPong:
		return nil
	case OpClose:
		return t.close(false)
	}
	return ErrNotControlFrame
}

func (t *TLVTransport) NextWriter(size uint32, typeID uint16, fn DataFrameWriterFunc) error {
	if size > t.maxFrameSize {
		return ErrFrameTooLarge
	}
	wHeader := MakeDataFrameHeader(size, typeID)
	err := WriteFrameHeader(wHeader, t.w)
	if err != nil {
		return err
	}
	if size == 0 {
		return flushWriter(t.w)
	}
	lw := &limitedWriter{W: t.w, N: int64(wHeader.PayloadSize())}
	err = fn(lw)
	if err != nil {
		return err
	}
	if lw.N != 0 {
		return ErrUncompletedWrite
	}
	return flushWriter(t.w)
}

func (t *TLVTransport) close(nextRead bool) error {
	err := WriteFrameHeader(closeFrame, t.w)
	if err != nil {
		return err
	}
	err = flushWriter(t.w)
	if err != nil {
		return err
	}
	if nextRead {
		_, err = ReadFrameHeader(t.r)
		if err != nil {
			return err
		}
	}
	return ErrClose
}

func (t *TLVTransport) Close() error {
	return t.close(true)
}

func (t *TLVTransport) Ping() error {
	err := WriteFrameHeader(pingFrame, t.w)
	if err != nil {
		return err
	}
	return flushWriter(t.w)
}
