package app

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// NewDataProvider returns DataProvider instance
func NewDataProvider(filePath string, skipHeader bool) *DataProvider {
	return &DataProvider{
		filePath:   filePath,
		skipHeader: skipHeader,
	}
}

// DataProvider is little abstraction over file
// Implements iterator pattern with streaming to keep memory usage low
type DataProvider struct {
	filePath   string
	err        error
	file       *os.File
	skipHeader bool
}

func (r *DataProvider) Open() error {
	file, err := os.Open(r.filePath)
	if err != nil {
		defer file.Close()
		return err
	}

	r.file = file
	return nil
}

func (r *DataProvider) Iter () <-chan []string {
	ch := make(chan []string)

	var src io.Reader = r.file
	ext := filepath.Ext(r.filePath)
	if ext == ".gz" {
		src, r.err = gzip.NewReader(src)
		if r.err != nil {
			close(ch)
			return ch
		}
	}

	reader := csv.NewReader(src)
	if r.skipHeader {
		if _, err := reader.Read(); err != nil {
			panic(err)
		}
	}

	go func () {
		for {
			row, err := reader.Read()
			if err == nil {
				ch <- row
				continue
			}

			if errors.Is(io.EOF, err) {
				err := r.file.Close()
				if err != nil {
					r.err = err
				}

				close(ch)
				break
			}

			r.err = err
		}
	} ()
	return ch
}

func (r *DataProvider) HasError() bool {
	return r.err != nil
}

func (r *DataProvider) Error() error {
	return r.err
}