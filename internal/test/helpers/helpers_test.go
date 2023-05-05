package helpers_test

import (
	"path/filepath"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestCreateTemp(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	if exists, err := tmp.Exists(); !exists {
		t.Errorf("Tmp file does not exist %s, err %s", tmp, err)
	}
}

func TestCreateTempInHome(t *testing.T) {
	tmp := helpers.CreateTempInHome(t)

	expected := filepath.Join("~", tmp.Base())
	if tmp.String() != expected {
		t.Errorf("CreateTempInHome() got %s, want %s", tmp, expected)
	}
}

func TestCreateTempDir(t *testing.T) {
	tmp := helpers.CreateTempDir(t)

	if exists, err := tmp.Exists(); !exists {
		t.Errorf("Tmp dir does not exist %s, err %s", tmp, err)
	}
}

func ensureExists(t *testing.T, p paths.Path, expected bool) {
	exists, err := p.Exists()
	if err != nil {
		t.Fatalf("Unexpected error checking existence path: %s, err: %s", p, err)
	}
	if exists != expected {
		t.Fatalf("File exist mismatch path: %s, want: %t, got: %t", p, expected, exists)
	}
}

// This test is here to avoid cyclical dependencies in the test
// because the path tests are in the same package as the code
// because they have to test private methods
func TestPathWrite(t *testing.T) {
	tests := []struct {
		parent      paths.Path
		description string
	}{
		{helpers.CreateTempDir(t), "Writes to temp file"},
		{helpers.CreateTempDirInHome(t), "Writes to home"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			p := paths.Path(filepath.Join(test.parent.String(), "subdir1/subdir2", "aaa.txt"))
			ensureExists(t, p, false)
			if err := p.Write([]byte("aaaa")); err != nil {
				t.Errorf("Write failed path: %s, err: %s", p.String(), err)
			}
			ensureExists(t, p, true)
		})
	}
}
