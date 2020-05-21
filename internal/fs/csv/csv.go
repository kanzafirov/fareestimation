package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const buffSize = 50

// NewReadChannel gets a file path as parameter, opens a csv file, reads it row by row,
// returns a new []string channel, where passes each row from the file
func NewReadChannel(path string) (<-chan []string, <-chan error, error) {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid file path: %v; err: %v ", path, err)
	}

	file, err := os.Open(fullpath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open file: %v; err: %v ", fullpath, err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	out := make(chan []string, buffSize)
	outerror := make(chan error)

	go func() {
		for {
			row, err := reader.Read()

			if err == io.EOF {
				break
			} else if err != nil {
				outerror <- fmt.Errorf("Failed to open file: %v; err: %v ", fullpath, err)
			}

			out <- row
		}

		file.Close()
		close(out)
	}()
	return out, outerror, nil
}

// NewWriteChannel gets a file name, input channel as parameter, passing string[] objects, representing csv rows
// and writes them in a file; uses sync.WaitGroup to notify when completed
func NewWriteChannel(path string, in <-chan []string) (<-chan int, <-chan error, error) {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid file path: %v; err: %v ", path, err)
	}

	dirpath := filepath.Dir(fullpath)
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not create path: %v; err: %v ", path, err)
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot create file; err: %v", err)
	}

	writer := csv.NewWriter(file)
	done := make(chan int)
	outerror := make(chan error)

	go func() {
		for row := range in {
			err := writer.Write(row)
			if err != nil {
				outerror <- fmt.Errorf("Cannot write row in file; err: %v", err)
			}
		}
		writer.Flush()
		file.Close()
		done <- 0
	}()
	return done, outerror, nil
}

// DeepCompare compares two files row by row without avoiding the load of the entire files in memory
func DeepCompare(file1, file2 string) (bool, error) {
	fullpathFile1, err := filepath.Abs(file1)
	if err != nil {
		log.Panicf("Invalid file1 path: %v; err: %v ", file1, err)
	}

	sf, err := os.Open(fullpathFile1)
	if err != nil {
		return false, fmt.Errorf("Failed to open file 1: %v; err: %v ", fullpathFile1, err)
	}

	fullpathFile2, err := filepath.Abs(file2)
	if err != nil {
		log.Panicf("Invalid file2 path: %v; err: %v ", file2, err)
	}

	df, err := os.Open(fullpathFile2)
	if err != nil {
		return false, fmt.Errorf("Failed to open file 2: %v; err: %v ", fullpathFile2, err)
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false, nil
		}
	}

	return true, nil
}
