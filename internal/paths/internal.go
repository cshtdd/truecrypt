package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Common path functions should be internal
type path string

func (p path) fullPath() string {
	if after, found := strings.CutPrefix(string(p), "~"); found {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Cannot read $HOME")
		}
		return filepath.Join(homeDir, after)
	}
	return string(p)
}

func (p path) exists() (bool, error) {
	if _, err := os.Stat(p.fullPath()); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (p path) delete() error {
	return os.RemoveAll(p.fullPath())
}

func (p path) base() string {
	return filepath.Base(p.fullPath())
}

func (p path) String() string {
	return string(p)
}
