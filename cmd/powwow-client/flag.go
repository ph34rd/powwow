package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type flags struct {
	*flag.FlagSet

	Addr         string
	PrintHelp    bool
	PrintVersion bool
}

func newFlagSet() *flags {
	f := &flags{}
	f.FlagSet = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	f.FlagSet.SetOutput(ioutil.Discard)
	f.StringVar(&f.Addr, "addr", "localhost:9999", "connection address")
	f.BoolVar(&f.PrintHelp, "help", false, "show help")
	f.BoolVar(&f.PrintVersion, "version", false, "print version information")
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
