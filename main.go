package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	message string
	logFile string
}

type FileLog struct {
	name       string
	lineOffset int64
}

func main() {
	filenames := []FileLog{
		{"testfile.log", 7},
		{"testfile_short.log", 1},
		{"testfile_empty.log", 1},
	}
	TICKINTERVAL := time.Second * 5
	results := make(chan Result, len(filenames))
	quit := make(chan string, 1)
	ERROR_REGEX, err := regexp.Compile("ERROR|WARN")
	if err != nil {
		panic("Invalid regex")
	}

	// Launching go routines for files
	for _, f := range filenames {
		// Doing this because of range variable f captured by func literal warning
		ff := f
		go func() {
			readFile(ff, results, TICKINTERVAL, ERROR_REGEX)
		}()
	}

	//Start timer
	go func() {
		<-time.After(time.Minute * 10)
		quit <- "End of timeout"
	}()

	//Wait for goroutines
	go func() {
		for i := 0; i < len(filenames); i++ {
			fmt.Println("======================")
			foundError := <-results
			fmt.Printf("This is the error found in file: %s\n", foundError.logFile)
			fmt.Println(foundError.message)
		}
		quit <- "All files completed"
	}()

	reason := <-quit
	fmt.Println(reason)

}

func readFile(flog FileLog, results chan Result, tickInterval time.Duration) {
	// Trying to read file
	f, err := os.Open(flog.name)
	if err != nil {
		results <- Result{logFile: flog.name, message: fmt.Sprintf("No error found because: %s", err)}
		return
	}
	defer f.Close()

	offset := getOffset(f, int(flog.lineOffset))
	nBytes := offset
	_, err = f.Seek(offset, 0)

	// Start ticker
	ticker := time.NewTicker(tickInterval)

	// Size of error snippet
	buf := make([]byte, 800)

	for range ticker.C {
		content, found, err := findError(f, buf, &nBytes)
		if err != nil {
			results <- Result{logFile: f.Name(), message: fmt.Sprintf("No error found because: %s", err)}
			return
		}
		if found {
			results <- Result{logFile: f.Name(), message: content}
			return
		}
	}

	ticker.Stop()
}

func findError(f *os.File, bufer []byte, offset *int64) (string, bool, error) {
	n2, err := f.Read(bufer)
	if err != nil {
		return "", false, err
	}

	content := string(bufer[:n2])

	if strings.Contains(content, "ERROR") {
		return content, true, nil
	}
	*offset += int64(n2)
	f.Seek(*offset, 0)
	return "", false, nil
}

func getOffset(reader io.Reader, lineNum int) int64 {
	scanner := bufio.NewScanner(reader)
	bytesRead := int64(0)
	for scanner.Scan() {
		currentSlice := scanner.Bytes()
		if bytes.Contains(currentSlice, []byte(strconv.Itoa(lineNum))) {
			return bytesRead
		}
		bytesRead += int64(len(currentSlice))
	}
	return bytesRead
}
