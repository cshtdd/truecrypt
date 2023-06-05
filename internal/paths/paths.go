package paths

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Match definitions

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

// Path definitions

type path string
type FilePath path // TODO: Use this more widely, or keep it internal maybe
type DirPath path
type ZipPath FilePath

type Path interface {
	FullPath() string
	Exists() (bool, error)
	Delete() error
	Base() string
	DirName() string
}

// Common path functions should be internal

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

func (p path) dirName() string {
	return filepath.Dir(p.fullPath())
}

func (p path) String() string {
	return string(p)
}

// FilePath functions

func CreateTempZipFile() (ZipPath, error) {
	if tmp, err := os.CreateTemp("", "tc_temp*.zip"); err != nil {
		return "", err
	} else {
		z := ZipPath(tmp.Name())
		if !z.IsValid() {
			return "", errors.New("invalid zip")
		}
		return z, nil
	}
}

func (p FilePath) FullPath() string {
	return path(p).fullPath()
}

func (p FilePath) Exists() (bool, error) {
	return path(p).exists()
}

func (p FilePath) Delete() error {
	return path(p).delete()
}

func (p FilePath) Base() string {
	return path(p).base()
}

func (p FilePath) DirName() string {
	return path(p).dirName()
}

func (p FilePath) String() string {
	return path(p).String()
}

func (p FilePath) Read() ([]byte, error) {
	return os.ReadFile(path(p).fullPath())
}

func (p FilePath) Write(data []byte) error {
	if err := p.CreateDir(); err != nil {
		return err
	}
	const userRwOthersR = 0644
	return os.WriteFile(path(p).fullPath(), data, userRwOthersR)
}

func (p FilePath) CreateDir() error {
	return DirPath(path(p).dirName()).Create()
}

func (p FilePath) Move(dest FilePath) error {
	if err := dest.CreateDir(); err != nil {
		return err
	}
	return os.Rename(path(p).fullPath(), path(dest).fullPath())
}

func (p FilePath) Matches(fileB FilePath) (MatchType, error) {
	existsA, err := path(p).exists()
	if err != nil {
		return Mismatch, err
	}

	existsB, err := path(fileB).exists()
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

// ZipPath functions

func (p ZipPath) FullPath() string {
	return path(p).fullPath()
}

func (p ZipPath) Exists() (bool, error) {
	return path(p).exists()
}

func (p ZipPath) Delete() error {
	return path(p).delete()
}

func (p ZipPath) Base() string {
	return path(p).base()
}

func (p ZipPath) DirName() string {
	return path(p).dirName()
}

func (p ZipPath) String() string {
	return path(p).String()
}

func (p ZipPath) IsValid() bool {
	return strings.HasSuffix(strings.ToLower(p.FullPath()), ".zip")
}

func (p ZipPath) Move(dest ZipPath) error {
	return FilePath(p).Move(FilePath(dest))
}

// DirPath functions

func (p DirPath) FullPath() string {
	return path(p).fullPath()
}

func (p DirPath) Exists() (bool, error) {
	return path(p).exists()
}

func (p DirPath) Delete() error {
	return path(p).delete()
}

func (p DirPath) Base() string {
	return path(p).base()
}

func (p DirPath) DirName() string {
	return path(p).dirName()
}

func (p DirPath) String() string {
	return path(p).String()
}

func (p DirPath) Create() error {
	fullPath := path(p).fullPath()
	const userRwOthersNone = 0700 // tmp files use this forcibly
	return os.MkdirAll(fullPath, userRwOthersNone)
}

func (p DirPath) Copy(dest DirPath) error {
	// wipe the dest beforehand
	if err := path(dest).delete(); err != nil {
		return err
	}

	// copy the files
	cmd := exec.Command("cp", "-r",
		fmt.Sprintf("%s/", path(p).fullPath()), path(dest).fullPath(),
	)
	_, err := cmd.Output()
	return err
}

func (p DirPath) Matches(dirB DirPath) MatchType {
	cmd := exec.Command("diff", "-q", "-r", path(p).fullPath(), path(dirB).fullPath())
	if _, err := cmd.Output(); err != nil {
		return Mismatch
	}
	return Match
}
