package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type filepaths []string

func (f *filepaths) String() string {
	return strings.Join(*f, ",")
}

func (f *filepaths) Set(value string) error {
	if value == "" {
		return errors.New("You need to provide a valid file")
	}
	*f = append(*f, value)
	return nil
}

type flags struct {
	filepaths   filepaths
	offset      int64
	interval    time.Duration
	limit       time.Duration
	pattern     string
	verbose     bool
	slack       bool
	output_file string
}

func (f *flags) parse() error {
	flag.Var(&f.filepaths, "file", "Log file path (required)")
	flag.Int64Var(&f.offset, "offset", f.offset, "Line number offset")
	flag.DurationVar(&f.interval, "interval", f.interval, "Interval at which the file will be checked")
	flag.DurationVar(&f.limit, "limit", f.limit, "By what time the program will be running")
	flag.StringVar(&f.pattern, "pattern", f.pattern, "Pattern to look for errors in file")
	flag.BoolVar(&f.verbose, "verbose", f.verbose, "Verbose mode, shows the errors in the standard output")
	flag.BoolVar(&f.slack, "slack", f.slack, "Send notification through slack")
	flag.StringVar(&f.output_file, "output_file", f.output_file, "Store found patterns in file (needs to be absolute path)")
	flag.Parse()

	if err := f.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}

func (f *flags) validate() error {
	if f.offset < 1 {
		return errors.New("The offset should be equal or greater than 1")
	}
	// TODO: Should there be a check for the minimum interval that we can have?
	return nil
}
