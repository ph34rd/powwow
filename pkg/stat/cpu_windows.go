package stat

import (
	"syscall"
	"time"
)

func ftToDuration(ft *syscall.Filetime) time.Duration {
	ns := ft.Nanoseconds()
	return time.Duration(ns)
}

// GetCPUUsage gathers current process times.
func GetCPUUsage() (CPUUsage, error) {
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return CPUUsage{}, err
	}
	var creationTime, exitTime, kernelTime, userTime syscall.Filetime
	err = syscall.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime)
	if err != nil {
		return CPUUsage{}, err
	}
	return CPUUsage{
		System: ftToDuration(&kernelTime),
		User:   ftToDuration(&userTime),
	}, nil
}
