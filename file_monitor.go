package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

type FileMonitor struct {
	filepath string
	offset   int64
	interval time.Duration
	regex    *regexp.Regexp
}

func (monitor *FileMonitor) Start(results chan Result) {
	// Trying to read file
	f, err := os.Open(monitor.filepath)
	if err != nil {
		results <- Result{filepath: monitor.filepath, message: fmt.Sprintf("No pattern found because: %s", err)}
		return
	}
	defer f.Close()

	offset := getOffset(f, int(monitor.offset))
	nBytes := offset

	// Move file offset just if file size is greater than offset
	fileInfo, err := f.Stat()
	if err != nil {
		results <- Result{filepath: monitor.filepath, message: fmt.Sprintf("No pattern found because: %s", err)}
		return
	}

	if fileInfo.Size() > offset {
		_, err = f.Seek(offset, 0)
	}

	// Start ticker
	ticker := time.NewTicker(monitor.interval)

	// Size of error snippet
	buf := make([]byte, 800)

	for range ticker.C {
		content, found, err := findError(f, buf, &nBytes, monitor.regex)
		if err != nil {
			results <- Result{filepath: monitor.filepath, message: fmt.Sprintf("No pattern found because: %s", err)}
			return
		}
		if found {
			results <- Result{filepath: monitor.filepath, message: content}
			return
		}
	}

	ticker.Stop()
}

func findError(f *os.File, bufer []byte, offset *int64, regex *regexp.Regexp) (string, bool, error) {
	n2, err := f.Read(bufer)
	// if we reach end to file wait for next ticker
	if errors.Is(err, io.EOF) {
		return "", false, nil
	}

	// if is another kind of error return the error
	if err != nil {
		return "", false, err
	}

	content := string(bufer[:n2])

	if regex.MatchString(content) {
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
