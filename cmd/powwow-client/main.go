package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ph34rd/powwow/pkg/session"
)

const (
	cmdName     = "powwow-client"
	dialTimeout = 15 * time.Second
)

var Version = "devel"

func printVersion() {
	fmt.Println(Version)
	os.Stdout.Sync()
}

func main() {
	time.Local = time.UTC

	fs := newFlagSet()
	err := fs.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "command line error: %v", err)
		os.Exit(1)
		return
	}

	if fs.PrintVersion {
		printVersion()
		return
	}
	if fs.PrintHelp {
		fs.PrintDefaults()
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancelFn := context.WithCancel(context.Background())

	d := net.Dialer{Timeout: dialTimeout, Deadline: time.Now().Add(dialTimeout)}
	conn, err := d.DialContext(ctx, "tcp", fs.Addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial err: %v\n", err)
		os.Exit(1)
		return
	}
	defer conn.Close()

	go func() {
		for range sigChan {
			cancelFn()
			conn.Close()
		}
	}()

	c := session.NewClientHandler(conn)
	wow, err := c.GetWoW(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get wow err: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println(wow)
}
