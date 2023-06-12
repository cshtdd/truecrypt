package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type IO struct {
	Reader io.Reader
	Writer io.Writer

	scanner          *bufio.Scanner
	scannerSensitive *bufio.Reader
}

func (io *IO) ReadLine() (read bool, line string, err error) {
	// use a reader singleton to avoid missing data
	if io.scanner == nil {
		io.scanner = bufio.NewScanner(io.Reader)
	}

	if read = io.scanner.Scan(); read {
		line = io.scanner.Text()
	}
	return read && len(line) > 0, line, io.scanner.Err()
}

// ReadSensitiveLine reads a line from the StdIn without echoing back.
// It maintains the interface and behavior as ReadLine
// Inspiration https://stackoverflow.com/a/48524979/247328
func (io *IO) ReadSensitiveLine() (read bool, line string, err error) {
	// use a reader singleton to avoid missing data
	if io.scannerSensitive == nil {
		io.scannerSensitive = bufio.NewReader(io.Reader)
	}

	// set the tty to raw and revert it back
	if err := io.runSttyWithArg("raw"); err != nil {
		return false, "", err
	}
	defer func() {
		if err := io.runSttyWithArg("-raw"); err != nil {
			panic(err)
		}
	}()

	// read the input
	var runes []rune
	for {
		switch c, _, err := io.scannerSensitive.ReadRune(); {
		case err != nil:
			return false, "", err
		case c == '\x03': // ctrl+c
			return false, "", nil
		case c == '\r' || c == '\n': // enter
			return len(runes) > 0, string(runes), nil
		default:
			runes = append(runes, c)
		}
	}
}

func (io *IO) WriteLine(s string) {
	_, _ = fmt.Fprintln(io.Writer, s)
}

func (io *IO) Pause() {
	io.WriteLine("Press any key to continue...")
	_, _, _ = io.ReadLine()
}

func (io *IO) runSttyWithArg(arg string) error {
	// There's no point in running this command against a non-shell
	// Moreover, stty will fail
	if io.Reader != os.Stdin {
		return nil
	}
	cmd := exec.Command("stty", arg)
	cmd.Stdin = io.Reader
	_, err := cmd.Output()
	return err
}
