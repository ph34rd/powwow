package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ph34rd/powwow/app"
)

const cmdName = "powwow-client"

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

	go func() {
		for range sigChan {
			cancelFn()
		}
	}()

	res, err := app.RunClient(ctx, fs.Addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client err: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println(res)
}
