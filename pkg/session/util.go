package session

import (
	"io"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	bb "github.com/valyala/bytebufferpool"

	"github.com/ph34rd/powwow/pkg/session/transport"
)

type linger interface {
	SetLinger(sec int) error
}

type marshalTo interface {
	proto.Message
	Size() int
	MarshalTo(dAtA []byte) (int, error)
}

func marshalHelper(w io.Writer, m marshalTo) error {
	buf := bb.Get()
	buf.Reset()
	defer bb.Put(buf)

	size := m.Size()
	if cap(buf.B) >= size {
		buf.B = buf.B[0:size]
		_, err := m.MarshalTo(buf.B)
		if err != nil {
			return err
		}
	} else {
		newBuf, err := proto.Marshal(m)
		if err != nil {
			return err
		}
		buf.Set(newBuf)
	}
	_, err := buf.WriteTo(w)
	return err
}

func frameWriterFunc(m marshalTo) transport.DataFrameWriterFunc {
	return func(w io.Writer) error {
		return marshalHelper(w, m)
	}
}

func unmarshalHelper(r io.Reader, m proto.Message) error {
	buf := bb.Get()
	defer bb.Put(buf)

	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buf.B, m)
}

type connDeadlineWrapper struct {
	conn    net.Conn
	timeout time.Duration
}

func newConnDeadlineWrapper(conn net.Conn, timeout time.Duration) *connDeadlineWrapper {
	return &connDeadlineWrapper{conn: conn, timeout: timeout}
}

func (c *connDeadlineWrapper) Read(p []byte) (n int, err error) {
	if c.timeout > 0 {
		err = c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return
		}
	}
	return c.conn.Read(p)
}

func (c *connDeadlineWrapper) Write(p []byte) (n int, err error) {
	if c.timeout > 0 {
		err = c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return
		}
	}
	return c.conn.Write(p)
}

func extractIP(conn net.Conn) string {
	addr := conn.RemoteAddr()
	tcp, ok := addr.(*net.TCPAddr)
	if ok {
		return tcp.IP.String()
	}
	return ""
}
