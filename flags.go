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
	filepaths    filepaths
	offset       int64
	interval     time.Duration
	limit        time.Duration
	pattern      string
	verbose      bool
	slack        bool
	output_file  string
	hide_finding bool
}

func (f *flags) parse() error {
	flag.Var(&f.filepaths, "file", "Log file path (required)")
	flag.Int64Var(&f.offset, "offset", f.offset, "Line number offset")
	flag.DurationVar(&f.interval, "interval", f.interval, "Interval at which the file will be checked")
	flag.DurationVar(&f.limit, "limit", f.limit, "By what time the program will be running")
	flag.StringVar(&f.pattern, "pattern", f.pattern, "Pattern to look for errors in file")
	flag.BoolVar(&f.verbose, "verbose", f.verbose, "Verbose mode, shows the errors in the standard output")
	flag.BoolVar(&f.slack, "slack", f.slack, "Send notification through slack")
	flag.BoolVar(&f.hide_finding, "hide_finding", f.hide_finding, "Hides the match and just notifies than an error was found")
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
	if !f.slack && f.output_file == "" && !f.verbose {
		return errors.New("You have to choose between notify through a slack message, store finding in file or show output in standard output")
	}
	// TODO: Should there be a check for the minimum interval that we can have?
	return nil
}
