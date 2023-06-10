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

func CreateTempZip(t *testing.T) paths.ZipPath {
	f := CreateTemp(t)
	return renameToZip(t, f)
}

func CreateTempInHome(t *testing.T) paths.FilePath {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}
	fullPath := createTempInDir(t, home)
	return paths.FilePath(filepath.Join("~", fullPath.Base()))
}

func CreateTempZipInHome(t *testing.T) paths.ZipPath {
	f := CreateTempInHome(t)
	return renameToZip(t, f)
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

func renameToZip(t *testing.T, f paths.FilePath) paths.ZipPath {
	t.Cleanup(func() {
		if err := f.Delete(); err != nil {
			panic(err)
		}
	})
	z := paths.ZipPath(fmt.Sprintf("%s.zip", f))
	if err := f.Move(paths.FilePath(z.FullPath())); err != nil {
		t.Fatalf("Error renaming temp file")
	}
	t.Cleanup(func() {
		if err := f.Delete(); err != nil {
			panic(err)
		}
	})
	return z
}

func createTempInDir(t *testing.T, dir string) paths.FilePath {
	tmp, err := os.CreateTemp(dir, "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
		return ""
	}
	p := paths.FilePath(tmp.Name())
	t.Cleanup(func() {
		p.Delete()
	})
	return p
}

func createTempDirInDir(t *testing.T, dir string) paths.DirPath {
	result, err := os.MkdirTemp(dir, "tc_test")
	if err != nil {
		t.Fatal("Cannot create temp dir")
		return ""
	}
	p := paths.DirPath(result)
	t.Cleanup(func() {
		p.Delete()
	})
	return p
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

func SetEnv(key string, value string, t *testing.T) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Unexpected error setting env: %s, err: %s", key, err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv(key); err != nil {
			panic(err)
		}
	})
}
