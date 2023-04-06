package internal

import "io"

type IO struct {
	Reader io.Reader
	Writer io.Writer
}
