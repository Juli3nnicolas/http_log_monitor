package reader

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// File describes a structure to read a file
type File struct {
	file    *os.File
	scanner *bufio.Scanner
	Parse   Parser
}

// Parser is a type of function used to interpret read files.
// All this Parsing system could be done better if instead of returning log.Info,
// []byte would be returned. However, it takes way more time to implement as it would require
// to code dedicated marshalers/unmarshalers
type Parser func(data []byte) (log.Info, error)

// CommonLogFormatParser is the default parser. It reads W3C httpd logs by default
func CommonLogFormatParser() Parser {
	return commonLogFormatParser
}

func commonLogFormatParser(data []byte) (log.Info, error) {
	return log.Parse(string(data))
}

// Open opens a file in read mode
func (r *File) Open(path ...interface{}) error {
	if len(path) != 1 {
		return fmt.Errorf("wrong argument number")
	}

	p, ok := path[0].(string)
	if !ok {
		return fmt.Errorf("invalid type - path must be a string")
	}

	f, err := os.Open(p)
	if err != nil {
		return err
	}

	r.file = f
	r.scanner = bufio.NewScanner(r.file)

	if r.Parse == nil {
		r.Parse = CommonLogFormatParser()
	}

	return nil
}

// Read reads a file content line by line
// Returns a nil slice when reaching EOF
func (r *File) Read() ([]log.Info, error) {
	if r.scanner.Scan() {
		v, err := r.Parse([]byte(r.scanner.Text()))
		if err != nil {
			return nil, err
		}
		return []log.Info{v}, nil
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

// Close closes the file opened with Open
func (r *File) Close() {
	if r.file != nil {
		r.file.Close()
		r.file = nil
	}
	r.scanner = nil
}
