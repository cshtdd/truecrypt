package helpers

import (
	"crypto/rand"
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

func getSamplePaths(dir paths.Path) []paths.Path {
	return []paths.Path{
		paths.Path(filepath.Join(dir.String(), "subdir1/subdir2", "aaa.txt")),
		paths.Path(filepath.Join(dir.String(), "subdir1/subdir2", "bbb.txt")),
		paths.Path(filepath.Join(dir.String(), "subdir3", "ccc.txt")),
	}
}

func CreateSampleNestedStructure(t *testing.T, dir paths.Path) {
	for _, p := range getSamplePaths(dir) {
		EnsureExists(t, p, false)
		var data []byte = make([]byte, 15)
		if _, err := rand.Read(data); err != nil {
			t.Errorf("Random data generation failed path: %s, err: %s", p.String(), err)
		}
		if err := p.Write(data); err != nil {
			t.Errorf("Write failed path: %s, err: %s", p.String(), err)
		}
		EnsureExists(t, p, true)
	}
}

func SampleDirectoriesMatch(t *testing.T, dirA paths.Path, dirB paths.Path) paths.MatchType {
	filesB := getSamplePaths(dirB)

	for index, fileA := range getSamplePaths(dirA) {
		fileB := filesB[index]
		switch m, err := fileA.Matches(fileB); {
		case err != nil:
			t.Fatalf("Unexpected error comparing files a: %s, b:%s, err: %s", fileA, fileB, err)
		case m != paths.Match:
			return m
		}
	}

	return paths.Match
}
