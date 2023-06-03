package paths

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Path string

func (p Path) FullPath() string {
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
	return os.ReadFile(p.FullPath())
}

func (p Path) CreateDir() error {
	fullPath := p.FullPath()
	const userRwOthersNone = 0700 // tmp files use this forcibly
	return os.MkdirAll(filepath.Dir(fullPath), userRwOthersNone)
}

func (p Path) Write(data []byte) error {
	if err := p.CreateDir(); err != nil {
		return err
	}
	const userRwOthersR = 0644
	return os.WriteFile(p.FullPath(), data, userRwOthersR)
}

func (p Path) MoveFile(dest Path) error {
	if err := dest.CreateDir(); err != nil {
		return err
	}
	return os.Rename(p.FullPath(), dest.FullPath())
}

func (p Path) CopyDir(dest Path) error {
	// wipe the dest beforehand
	if err := dest.Delete(); err != nil {
		return err
	}

	// copy the files
	cmd := exec.Command("cp", "-r",
		fmt.Sprintf("%s/", p.FullPath()), dest.FullPath(),
	)
	_, err := cmd.Output()
	return err
}

func (p Path) Exists() (bool, error) {
	if _, err := os.Stat(p.FullPath()); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (p Path) Delete() error {
	return os.RemoveAll(p.FullPath())
}

func (p Path) Base() string {
	return filepath.Base(p.FullPath())
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

func (p Path) MatchesDir(dirB Path) MatchType {
	cmd := exec.Command("diff", "-q", "-r", p.FullPath(), dirB.FullPath())
	if _, err := cmd.Output(); err != nil {
		return Mismatch
	}
	return Match
}

func (p Path) MatchesFile(fileB Path) (MatchType, error) {
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

func CreateTempFile() (Path, error) {
	//TODO: remove duplication with the test helpers temp methods
	if tmp, err := os.CreateTemp("", "tc_temp"); err != nil {
		return "", err
	} else {
		return Path(tmp.Name()), nil
	}
}
