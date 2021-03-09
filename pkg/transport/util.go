package transport

import (
	"io"
	"io/ioutil"
)

type flusher interface {
	Flush() error
}

type nullReader struct{}

func (nullReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

func discardReader(r *io.LimitedReader) error {
	if r.N == 0 {
		return nil
	}
	_, err := io.Copy(ioutil.Discard, r)
	if err == nil {
		return err
	}
	if r.N != 0 {
		return ErrUnexpectedEOF
	}
	return nil
}

// A limitedWriter writes to W but limits the amount of data wrote to just N bytes.
type limitedWriter struct {
	W io.Writer // underlying reader
	N int64     // max bytes remaining
}

// Write implements io.Writer.
func (lw *limitedWriter) Write(p []byte) (n int, err error) {
	if int64(len(p)) > lw.N {
		return 0, ErrThresholdExceeded
	}
	n, err = lw.W.Write(p)
	if err == nil {
		lw.N -= int64(n)
		if n < 0 {
			return 0, ErrThresholdExceeded
		}
	}
	return
}

func flushWriter(w io.Writer) error {
	f, ok := w.(flusher)
	if ok {
		return f.Flush()
	}
	return nil
}

func (lw *limitedWriter) Flush() error {
	return flushWriter(lw.W)
}
