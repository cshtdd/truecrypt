package paths_test

import (
	"fmt"
	"os"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/test/helpers"
	"testing"
)

func TestFullPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}

	tests := []struct {
		path        string
		expected    string
		description string
	}{
		{"aaa.txt", "aaa.txt", "Does not modify plain files"},
		{"/aaa/aaa.txt", "/aaa/aaa.txt", "Does not modify full paths"},
		{"~", home, "Expands tilde"},
		{"~/aaa.txt", fmt.Sprintf("%s/aaa.txt", home), "Expands tilde for files"},
		{"~/bbb/aaa.txt", fmt.Sprintf("%s/bbb/aaa.txt", home), "Expands tilde for files in hierarchy"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := paths.Path(test.path).FullPath()
			if actual != test.expected {
				t.Errorf("paths.Expand(%s) = %s, want %s", test.path, actual, test.expected)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		parent      paths.Path
		description string
	}{
		{helpers.CreateTempDir(t), "Writes to temp file"},
		{helpers.CreateTempDirInHome(t), "Writes to home"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			helpers.CreateSampleNestedStructure(t, test.parent)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		fillSubdirs bool
		path        paths.Path
		description string
	}{
		{false, helpers.CreateTemp(t), "Deletes file"},
		{false, helpers.CreateTempInHome(t), "Deletes file in home"},
		{false, helpers.CreateTempDir(t), "Deletes empty dir"},
		{false, helpers.CreateTempDirInHome(t), "Deletes empty dir in home"},
		{true, helpers.CreateTempDir(t), "Deletes dir"},
		{true, helpers.CreateTempDirInHome(t), "Deletes dir in home"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.fillSubdirs {
				helpers.CreateSampleNestedStructure(t, test.path)
			}
			helpers.EnsureExists(t, test.path, true)

			if err := test.path.Delete(); err != nil {
				t.Errorf("Unexpected error deleting path: %s, err: %s", test.path, err)
			}

			helpers.EnsureExists(t, test.path, false)
		})
	}
}

func TestMoveFile(t *testing.T) {
	existingFile := helpers.CreateTemp(t)
	helpers.WriteRandomData(t, existingFile, helpers.GenerateRandomData(t))

	tests := []struct {
		src         paths.Path
		dest        paths.Path
		shouldFail  bool
		description string
	}{
		{
			"not_found.txt",
			"~/will_not_exist",
			true,
			"Cannot move non existent",
		},
		{
			helpers.CreateTemp(t),
			helpers.CreateTempInHome(t),
			false,
			"Expand dest paths",
		},
		{
			helpers.CreateTempInHome(t),
			helpers.CreateTemp(t),
			false,
			"Expand src paths",
		},
		{
			helpers.CreateTempInHome(t),
			paths.Path(fmt.Sprintf("%s/dir1/dir2/dir3/aaa.txt", helpers.CreateTempDir(t))),
			false,
			"Creates directory structure",
		},
		{
			helpers.CreateTemp(t),
			existingFile,
			false,
			"Overwrites existing files",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			seedData := helpers.GenerateRandomData(t)

			switch exists, err := test.src.Exists(); {
			case err != nil:
				t.Fatalf("Error checking existence of %s, err %s", test.src, err)
			case exists:
				helpers.WriteRandomData(t, test.src, seedData)
			}

			if err := test.src.MoveFile(test.dest); err != nil {
				if test.shouldFail {
					return
				}

				t.Fatalf(
					"Unexpected error moving src: %s to dest: %s, err: %s",
					test.src, test.dest, err,
				)
			}

			switch exists, err := test.src.Exists(); {
			case err != nil:
				t.Fatalf("Error checking existence of %s, err %s", test.src, err)
			case exists:
				t.Fatalf("Source should not exist after move")
			}

			switch exists, err := test.dest.Exists(); {
			case err != nil:
				t.Fatalf("Error checking existence of %s, err %s", test.dest, err)
			case !exists:
				t.Fatalf("Dest should exist after move")
			}

			control := helpers.CreateTemp(t)
			helpers.WriteRandomData(t, control, seedData)
			switch m, err := control.Matches(test.dest); {
			case err != nil:
				t.Fatalf("Unexpected error comparing files a: %s, b:%s, err: %s", control, test.src, err)
			case m != paths.Match:
				t.Fatalf("Dest path doesn't match control")
			}
		})
	}
}
