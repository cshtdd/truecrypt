package paths

import (
	"os"
	"path/filepath"
	"strings"
)

type Path string

func (p Path) expand() string {
	if after, found := strings.CutPrefix(string(p), "~"); found {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Cannot read $HOME")
		}
		return filepath.Join(homeDir, after)
	}
	return string(p)
}

func (p Path) Read() ([]byte, error) {
	return os.ReadFile(p.expand())
}

func (p Path) Write(data []byte) error {
	fullPath := p.expand()
	const userRwOthersNone = 0700 // tmp files use this forcibly
	os.MkdirAll(filepath.Dir(fullPath), userRwOthersNone)
	const userRwOthersR = 0644
	return os.WriteFile(fullPath, data, userRwOthersR)
}

func (p Path) Exists() (bool, error) {
	if _, err := os.Stat(p.expand()); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (p Path) Base() string {
	return filepath.Base(p.expand())
}

func (p Path) String() string {
	return string(p)
}
