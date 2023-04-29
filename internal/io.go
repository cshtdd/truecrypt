package internal

import (
	"bufio"
	"io"
)

type IO struct {
	Reader io.Reader
	Writer io.Writer
}

type LineReader bufio.Scanner

func (io *IO) NewLineReader() *LineReader {
	return (*LineReader)(bufio.NewScanner(io.Reader))
}

func (s *LineReader) Read() (read bool, line string, err error) {
	scanner := (*bufio.Scanner)(s)
	if read = scanner.Scan(); read {
		line = scanner.Text()
	}
	return read, line, scanner.Err()
}
