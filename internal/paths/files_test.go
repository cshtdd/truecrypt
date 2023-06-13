package paths_test

import (
	"fmt"
	"github.com/cshtdd/truecrypt/internal/paths"
	"github.com/cshtdd/truecrypt/internal/test/helpers"
	"os"
	"testing"
)

func TestFilePath_FullPath(t *testing.T) {
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
		{"~/aaa.txt", fmt.Sprintf("%s/aaa.txt", home), "Expands tilde for files"},
		{"~/bbb/aaa.txt", fmt.Sprintf("%s/bbb/aaa.txt", home), "Expands tilde for files in hierarchy"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := paths.FilePath(test.path).FullPath()
			if actual != test.expected {
				t.Errorf("paths.Expand(%s) = %s, want %s", test.path, actual, test.expected)
			}
		})
	}
}

func TestFilePath_CreateDir(t *testing.T) {
	t1 := helpers.CreateTempDir(t)
	t2 := helpers.CreateTempDirInHome(t)

	tests := []struct {
		p           paths.FilePath
		target      paths.DirPath
		description string
	}{
		{
			paths.FilePath(fmt.Sprintf("%s/aaa.txt", t1)),
			t1,
			"Does not create existing dir",
		},
		{
			paths.FilePath(fmt.Sprintf("%s/bbb/ccc/aaa.txt", t1)),
			paths.DirPath(fmt.Sprintf("%s/bbb/ccc/", t1)),
			"Creates dir structure"},
		{
			paths.FilePath(fmt.Sprintf("%s/bbb/ccc/ddd/eee/", t1)),
			paths.DirPath(fmt.Sprintf("%s/bbb/ccc/ddd/eee/", t1)),
			"Creates dir path"},
		{
			paths.FilePath(fmt.Sprintf("%s/aaa.txt", t2)),
			t2,
			"Does not create existing home dir",
		},
		{
			paths.FilePath(fmt.Sprintf("%s/bbb/ccc/aaa.txt", t1)),
			paths.DirPath(fmt.Sprintf("%s/bbb/ccc/", t1)),
			"Creates home dir structure"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if err := test.p.CreateDir(); err != nil {
				t.Fatalf("Unexpected error creating p: %s err: %s", test.p, err)
			}

			if exists, err := test.target.Exists(); !exists || err != nil {
				t.Fatalf("Target path does not exist p: %s exists: %t err: %s", test.target, exists, err)
			}
		})
	}
}

func TestFilePath_Write(t *testing.T) {
	tests := []struct {
		parent      paths.DirPath
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

func TestFilePath_Delete(t *testing.T) {
	tests := []struct {
		path        paths.FilePath
		description string
	}{
		{helpers.CreateTemp(t), "Deletes file"},
		{helpers.CreateTempInHome(t), "Deletes file in home"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			helpers.EnsureExists(t, test.path, true)

			if err := test.path.Delete(); err != nil {
				t.Errorf("Unexpected error deleting path: %s, err: %s", test.path, err)
			}

			helpers.EnsureExists(t, test.path, false)
		})
	}
}

func TestFilePath_Move(t *testing.T) {
	existingFile := helpers.CreateTemp(t)
	helpers.WriteRandomData(t, existingFile, helpers.GenerateRandomData(t))

	tests := []struct {
		src         paths.FilePath
		dest        paths.FilePath
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
			paths.FilePath(fmt.Sprintf("%s/dir1/dir2/dir3/aaa.txt", helpers.CreateTempDir(t))),
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

			if err := test.src.Move(test.dest); err != nil {
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

func TestZipPath_IsValid(t *testing.T) {
	tests := []struct {
		path        paths.ZipPath
		valid       bool
		description string
	}{
		{"/tmp/aaa.zip", true, "full zip file"},
		{"tmp/aaa.zip", true, "relative zip file"},
		{"~/tmp/aaa.zip", true, "home zip file"},
		{"~/tmp/aaa.ZIP", true, "uppercase zip file"},
		{"~/tmp/aaazip", false, "not a zip file"},
		{"aaa.txt", false, "different extension"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if v := test.path.IsValid(); v != test.valid {
				t.Errorf("IsValid(%s) want: %t got: %t", test.path, test.valid, v)
			}
		})
	}
}
