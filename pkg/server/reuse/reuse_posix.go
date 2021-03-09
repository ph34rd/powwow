//go:build !plan9 && !windows && !wasm
// +build !plan9,!windows,!wasm

package reuse

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func Control(_, _ string, c syscall.RawConn) (err error) {
	return c.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err == nil {
			return
		}
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if err != nil {
			return
		}
	})
}
