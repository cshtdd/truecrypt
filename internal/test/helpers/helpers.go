package helpers

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
)

func CreateTemp(t *testing.T) paths.FilePath {
	return createTempInDir(t, "")
}

func CreateTempInHome(t *testing.T) paths.FilePath {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}
	fullPath := createTempInDir(t, home)
	return paths.FilePath(filepath.Join("~", fullPath.Base()))
}

func CreateTempDir(t *testing.T) paths.DirPath {
	return createTempDirInDir(t, "")
}

func CreateTempDirInHome(t *testing.T) paths.DirPath {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}
	fullPath := createTempDirInDir(t, home)
	return paths.DirPath(filepath.Join("~", fullPath.Base()))
}

func createTempInDir(t *testing.T, dir string) paths.FilePath {
	tmp, err := os.CreateTemp(dir, "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
		return ""
	}
	t.Cleanup(func() {
		os.Remove(tmp.Name())
	})
	return paths.FilePath(tmp.Name())
}

func createTempDirInDir(t *testing.T, dir string) paths.DirPath {
	result, err := os.MkdirTemp(dir, "tc_test")
	if err != nil {
		t.Fatal("Cannot create temp dir")
		return ""
	}
	t.Cleanup(func() {
		os.RemoveAll(result)
	})
	return paths.DirPath(result)
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

func getSamplePaths(t *testing.T, dir paths.DirPath) []paths.FilePath {
	var result []paths.FilePath
	count := int(GenerateRandomData(t)[0]%5) + 1
	for i := 0; i < count; i++ {
		result = append(result, paths.FilePath(filepath.Join(
			dir.String(),
			"subdir1/subdir2",
			fmt.Sprintf("aaa%d.txt", i),
		)))
		result = append(result, paths.FilePath(filepath.Join(
			dir.String(),
			"subdir3",
			fmt.Sprintf("ccc%d.txt", i),
		)))
	}
	return result
}

func CreateSampleNestedStructure(t *testing.T, dir paths.DirPath) {
	for _, p := range getSamplePaths(t, dir) {
		EnsureExists(t, p, false)
		data := GenerateRandomData(t)
		WriteRandomData(t, p, data)
		EnsureExists(t, p, true)
	}
}

func WriteRandomData(t *testing.T, p paths.FilePath, data []byte) {
	if err := p.Write(data); err != nil {
		t.Errorf("Write failed path: %s, err: %s", p.String(), err)
	}
}

func GenerateRandomData(t *testing.T) []byte {
	var data = make([]byte, 15)
	if _, err := rand.Read(data); err != nil {
		t.Errorf("Random data generation failed err: %s", err)
	}
	return data
}
