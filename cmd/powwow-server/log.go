package main

import (
	"os"

	"go.uber.org/zap"
)

func defaultProcessFields(name string) []zap.Field {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return []zap.Field{
		zap.Int("pid", os.Getpid()),
		zap.String("host", host),
		zap.String("app", name),
	}
}
