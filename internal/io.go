package internal

import (
	"fmt"
	"io"
)

type IO struct {
	Reader io.Reader
	Writer io.Writer
}

func (w *IO) Writeln(s string) (int, error) {
	return w.Writer.Write([]byte(fmt.Sprintln(s)))
}
