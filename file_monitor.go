package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type FileProtector struct {
	filepath string
	mu       sync.Mutex
}

type FileMonitor struct {
	filepath  string
	offset    int64
	interval  time.Duration
	regex     *regexp.Regexp
	storeFile string
}

type Result struct {
	filepath string
	error    error
	found    bool
	content  string
}

func (monitor *FileMonitor) Start(results chan Result) {
	// Trying to read file
	f, err := os.Open(monitor.filepath)
	if err != nil {
		results <- Result{filepath: monitor.filepath, error: err}
		return
	}
	defer f.Close()

	record := make(map[string]int64)
	var offset int64
	fileProtector := &FileProtector{filepath: monitor.storeFile}
	initialOffset := fileProtector.getOffsetFromFile(monitor.filepath, record)
	fmt.Printf("This is the initial offset: %d\n", initialOffset)

	// We have a biggest offset stored in json file so use that one
	if initialOffset <= monitor.offset {
		offset = getOffset(f, int(monitor.offset))
	}
	fmt.Printf("offset: %d\n", offset)
	// nBytes := offset

	// Move file offset just if file size is greater than offset
	fileInfo, err := f.Stat()
	if err != nil {
		results <- Result{filepath: monitor.filepath, error: err}
		return
	}

	if fileInfo.Size() > offset {
		_, err = f.Seek(offset, 0)
	}

	// Size of error snippet
	buf := make([]byte, 800)

	content, found, err := findError(f, buf, &offset, monitor.regex, fileProtector, record, monitor.filepath)
	if err != nil {
		results <- Result{filepath: monitor.filepath, error: err}
		return
	}
	if found {
		results <- Result{filepath: monitor.filepath, found: true, content: content}
		return
	}
	if !found && err == nil {
		results <- Result{filepath: monitor.filepath, found: false, error: nil}
	}

}

func findError(f *os.File, bufer []byte, offset *int64, regex *regexp.Regexp, fileProtector *FileProtector, record map[string]int64, filepath string) (string, bool, error) {
	fmt.Println("Calling find error1")
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
		fmt.Println("FOUND ERROR!!")
		return content, true, nil
	}
	*offset += int64(n2)
	fmt.Printf("offset: %d\n", *offset)
	record[filepath] = *offset
	fmt.Printf("%q", record)
	fileProtector.storeOffsetInFile(record)
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

func (p *FileProtector) getOffsetFromFile(filepath string, record map[string]int64) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	f, err := os.OpenFile(p.filepath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("Error while trying to get offset: %s - %s\n", p.filepath, err)
		panic("Some error happened")
	}
	defer f.Close()

	byteValue, _ := os.ReadFile(p.filepath)
	fmt.Printf("Reading from file: %s -> %s", p.filepath, string(byteValue))
	// myMap := make(map[string]int)
	json.Unmarshal(byteValue, &record)
	fmt.Println(record)
	fmt.Println(record["offsets.json"])
	fmt.Printf("This is the filepath: %s\n", filepath)

	if value, ok := record[filepath]; ok {
		return value
	}
	return 1
}

func (p *FileProtector) storeOffsetInFile(value map[string]int64) {
	fmt.Println("Calling store in file")
	fmt.Println("What is going to be writen:")
	fmt.Println(value)
	p.mu.Lock()
	defer p.mu.Unlock()
	f, err := os.Open(p.filepath)
	if err != nil {
		panic("Error while trying to get offset")
	}
	defer f.Close()

	bytes, err := json.Marshal(value)
	if err != nil {
		panic("Error when trying to write offset")
	}
	err = ioutil.WriteFile(p.filepath, bytes, 06444)
	if err != nil {
		panic("Error when trying to write offset")
	}
}
