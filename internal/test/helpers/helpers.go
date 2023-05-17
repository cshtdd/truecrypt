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

func EnsureExists(t *testing.T, p paths.Path, expected bool) {
	exists, err := p.Exists()
	if err != nil {
		t.Fatalf("Unexpected error checking existence path: %s, err: %s", p, err)
	}
	if exists != expected {
		t.Fatalf("File exist mismatch path: %s, want: %t, got: %t", p, expected, exists)
	}
}

func CreateSampleNestedStructure(t *testing.T, dir paths.Path) {
	files := []paths.Path{
		paths.Path(filepath.Join(dir.String(), "subdir1/subdir2", "aaa.txt")),
		paths.Path(filepath.Join(dir.String(), "subdir1/subdir2", "bbb.txt")),
		paths.Path(filepath.Join(dir.String(), "subdir3", "ccc.txt")),
	}
	for _, p := range files {
		EnsureExists(t, p, false)
		if err := p.Write([]byte("aaaa")); err != nil {
			t.Errorf("Write failed path: %s, err: %s", p.String(), err)
		}
		EnsureExists(t, p, true)
	}
}
