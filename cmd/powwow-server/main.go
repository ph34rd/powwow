package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ph34rd/powwow/pkg/server"
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

	loggerConfig := zap.NewProductionConfig()
	if fs.Dev {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	logger, err := loggerConfig.Build(
		zap.Fields(defaultProcessFields(cmdName)...),
		zap.AddCaller(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init error: %v", err)
		os.Exit(1)
		return
	}
	defer logger.Sync()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	srv, err := server.NewServer(logger, fs.Bind)
	if err != nil {
		logger.Fatal("server init error", zap.Error(err))
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
		logger.Fatal("server error", zap.Error(err))
		os.Exit(1)
		return
	}
}
