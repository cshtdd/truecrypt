package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
)

func CreateTemp(t *testing.T) paths.Path {
	return createTempInDir(t, "")
}

func CreateTempInHome(t *testing.T) paths.Path {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}
	fullPath := createTempInDir(t, home)
	return paths.Path(filepath.Join("~", filepath.Base(fullPath.String())))
}

func CreateTempDir(t *testing.T) paths.Path {
	return createTempDirInDir(t, "")
}

func CreateTempDirInHome(t *testing.T) paths.Path {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}
	fullPath := createTempDirInDir(t, home)
	return paths.Path(filepath.Join("~", filepath.Base(fullPath.String())))
}

func createTempInDir(t *testing.T, dir string) paths.Path {
	tmp, err := os.CreateTemp(dir, "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
		return ""
	}
	t.Cleanup(func() {
		os.Remove(tmp.Name())
	})
	return paths.Path(tmp.Name())
}

func createTempDirInDir(t *testing.T, dir string) paths.Path {
	result, err := os.MkdirTemp(dir, "tc_test")
	if err != nil {
		t.Fatal("Cannot create temp dir")
		return ""
	}
	t.Cleanup(func() {
		os.RemoveAll(result)
	})
	return paths.Path(result)
}
