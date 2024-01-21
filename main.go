package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

type Result struct {
	message  string
	filepath string
}

type FileLog struct {
	name       string
	lineOffset int64
}

func main() {
	f := &flags{
		offset:   1,
		interval: time.Second * 5,
		limit:    time.Hour * 10,
		pattern:  "ERROR|WARN",
	}

	if err := f.parse(); err != nil {
		os.Exit(1)
	}
	ERROR_REGEX, err := regexp.Compile(f.pattern)
	if err != nil {
		panic("Invalid regex")
	}

	results := make(chan Result, len(f.filepaths))
	quit := make(chan string, 1)

	// Launching go routines for files
	for _, n := range f.filepaths {
		// Doing this because of range variable f captured by func literal warning
		ff := n
		go func() {
			monitor := &FileMonitor{offset: f.offset,
				interval: f.interval,
				regex:    ERROR_REGEX,
				filepath: ff}
			monitor.Start(results)
		}()
	}

	// Start timer
	go func() {
		<-time.After(f.limit)
		quit <- "End of timeout"
	}()

	// Wait for goroutines
	go func() {
		for i := 0; i < len(f.filepaths); i++ {
			fmt.Println("======================")
			foundError := <-results
			fmt.Printf("This is the error found in file: %s\n", foundError.filepath)
			fmt.Println(foundError.message)
		}
		quit <- "All files completed"
	}()

	reason := <-quit
	fmt.Println(reason)

}
