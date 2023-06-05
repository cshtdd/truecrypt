package paths

import (
	"fmt"
	"os"
	"os/exec"
)

type DirPath path

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
