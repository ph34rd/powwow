package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ph34rd/powwow/app"
	"github.com/ph34rd/powwow/pkg/logger"
)

const cmdName = "powwow-server"

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

	lg, err := logger.NewLogger(cmdName, fs.Dev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init error: %v", err)
		os.Exit(1)
		return
	}
	defer lg.Sync()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	srv, err := app.NewServer(lg, fs.Bind)
	if err != nil {
		lg.Fatal("server init error", zap.Error(err))
		os.Exit(1)
		return
	}

	go func() {
		for range sigChan {
			srv.Stop()
		}
	}()

	err = srv.Run()
	if err != nil {
		lg.Fatal("server error", zap.Error(err))
		os.Exit(1)
		return
	}
}
