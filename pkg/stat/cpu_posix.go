//go:build !windows
// +build !windows

package stat

import (
	"syscall"
	"time"
)

func tvToDuration(tv syscall.Timeval) time.Duration {
	return time.Duration(tv.Nano()) * time.Nanosecond
}

// GetCPUUsage gathers current process times.
func GetCPUUsage() (CPUUsage, error) {
	ru := syscall.Rusage{}
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	if err != nil {
		return CPUUsage{}, err
	}
	return CPUUsage{
		System: tvToDuration(ru.Stime),
		User:   tvToDuration(ru.Utime),
	}, nil
}
