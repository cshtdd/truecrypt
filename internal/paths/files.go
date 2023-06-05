package paths

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type FilePath path
type ZipPath FilePath

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

func (p FilePath) Dir() DirPath {
	return DirPath(filepath.Dir(p.FullPath()))
}

func (p FilePath) CreateDir() error {
	return p.Dir().Create()
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

func (p ZipPath) String() string {
	return path(p).String()
}

func (p ZipPath) IsValid() bool {
	return strings.HasSuffix(strings.ToLower(p.FullPath()), ".zip")
}

func (p ZipPath) Move(dest ZipPath) error {
	return FilePath(p).Move(FilePath(dest))
}
