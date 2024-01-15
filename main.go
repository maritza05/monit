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

func main() {
	f, err := os.Open("./testfile.log")
	if err != nil {
		panic("Error while trying to read file")
	}
	defer f.Close()
	offset := getOffset(f, 7)
	nBytes := offset
	fmt.Printf("++++This is the offtset: %d", offset)
	_, err = f.Seek(offset, 0)
	ticker := time.NewTicker(time.Second * 10)
	buf := make([]byte, 800)
	done := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				fmt.Println("We got to done!!!!")
				done <- true
				return
			case <-ticker.C:
				findError(f, buf, done, &nBytes)
			}
		}
	}()

	<-done
	fmt.Println("Stopping program!!!")
	ticker.Stop()
	close(done)
}

func findError(f *os.File, bufer []byte, done chan bool, offset *int64) {
	fmt.Println("Calling function!")
	n2, err := f.Read(bufer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Some error happened!")
			done <- true
			return
		}
	}

	content := string(bufer[:n2])

	fmt.Printf("++ %s", content)
	if strings.Contains(content, "ERROR") {
		fmt.Printf("Found error at!! %s", content)
		done <- true
		return
	}
	*offset += int64(n2)
	f.Seek(*offset, 0)
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
