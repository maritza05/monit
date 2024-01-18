package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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
	filenames := []string{
		"testfile.log",
		"testfile_short.log",
		"testfile_empty.log",
	}
	results := make(chan Result, len(filenames))
	quit := make(chan string)

	// Launching go routines for files
	for _, f := range filenames {
		// Doing this because of range variable f captured by func literal warning
		ff := f
		fmt.Printf("Type %T, value: %s\n", ff, ff)
		go func() {
			fmt.Printf("--> Launching go routine for: %s\n", ff)
			readFile(ff, results)
		}()
	}

	// Start timer
	go func() {
		<-time.After(time.Minute * 10)
		quit <- "End of timeout"
	}()

	// Wait for goroutines
	go func() {
		for result := range results {
			fmt.Println("======================")
			fmt.Printf("This is the error found in file: %s\n", result.logFile)
			fmt.Println(result.message)
		}
		quit <- "All files completed"
	}()

	reason := <-quit
	fmt.Println(reason)

}

func readFile(filename string, results chan Result) {
	fmt.Printf("Trying to read: %s\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		panic("Error while trying to read file")
	}
	defer f.Close()
	offset := getOffset(f, 1)
	nBytes := offset
	fmt.Printf("++++This is the offtset: %d", offset)
	_, err = f.Seek(offset, 0)
	ticker := time.NewTicker(time.Second * 5)
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
	fmt.Println("Calling function!")
	n2, err := f.Read(bufer)
	if err != nil {
		fmt.Println("Some error happened!")
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
