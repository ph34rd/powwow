package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type flags struct {
	*flag.FlagSet

	Bind         string
	PrintHelp    bool
	PrintVersion bool
	Dev          bool
}

func newFlagSet() *flags {
	f := &flags{}
	f.FlagSet = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	f.FlagSet.SetOutput(io.Discard)
	f.StringVar(&f.Bind, "bind", ":9999", "bind port")
	f.BoolVar(&f.PrintHelp, "help", false, "show help")
	f.BoolVar(&f.PrintVersion, "version", false, "print version information")
	f.BoolVar(&f.Dev, "dev", false, "enable development logging")
	return f
}

func (f *flags) Parse() error {
	err := f.FlagSet.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		return errors.New("flag provided but not defined: -h")
	}
	return err
}

func (f *flags) PrintDefaults() {
	fmt.Fprintf(os.Stdout, "usage: %s [options]\noptions:\n", cmdName)
	f.FlagSet.SetOutput(os.Stdout)
	f.FlagSet.PrintDefaults()
	os.Stdout.Sync()
}
