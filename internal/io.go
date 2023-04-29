package internal

import (
	"bufio"
	"fmt"
	"io"
)

type IO struct {
	Reader io.Reader
	Writer io.Writer

	scanner *bufio.Scanner
}

func (io *IO) ReadLine() (read bool, line string, err error) {
	if io.scanner == nil {
		io.scanner = bufio.NewScanner(io.Reader)
	}

	if read = io.scanner.Scan(); read {
		line = io.scanner.Text()
	}
	return read && len(line) > 0, line, io.scanner.Err()
}

func (io *IO) WriteLine(s string) {
	fmt.Fprintln(io.Writer, s)
}
