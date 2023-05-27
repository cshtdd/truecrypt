package paths

import (
	"bytes"
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

func (p Path) Delete() error {
	return os.RemoveAll(p.expand())
}

func (p Path) Base() string {
	return filepath.Base(p.expand())
}

func (p Path) String() string {
	return string(p)
}

type MatchType int

const (
	Mismatch MatchType = iota
	Match
)

func (m MatchType) String() string {
	switch m {
	case Match:
		return "Match"
	default:
		return "Mismatch"
	}
}

func (p Path) Matches(fileB Path) (MatchType, error) {
	existsA, err := p.Exists()
	if err != nil {
		return Mismatch, err
	}

	existsB, err := fileB.Exists()
	if err != nil {
		return Mismatch, err
	}

	if existsA != existsB { // mismatch
		return Mismatch, nil
	}

	if existsA && existsB {
		bytesA, err := p.Read()
		if err != nil {
			return Mismatch, err
		}

		bytesB, err := fileB.Read()
		if err != nil {
			return Mismatch, err
		}

		if !bytes.Equal(bytesA, bytesB) { // mismatch
			return Mismatch, nil
		}
	}

	return Match, nil
}
