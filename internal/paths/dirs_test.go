package paths_test

import (
	"fmt"
	"os"
	"path/filepath"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/test/helpers"
	"testing"
)

func TestDirPath_Delete(t *testing.T) {
	tests := []struct {
		fillSubdirs bool
		path        paths.DirPath
		description string
	}{
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

func TestDirPath_FullPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}

	tests := []struct {
		path        string
		expected    string
		description string
	}{
		{"/aaa/aaa", "/aaa/aaa", "Does not modify full paths"},
		{"~", home, "Expands tilde"},
		{"~/bbb/aaa", fmt.Sprintf("%s/bbb/aaa", home), "Expands tilde for files in hierarchy"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := paths.DirPath(test.path).FullPath()
			if actual != test.expected {
				t.Errorf("paths.Expand(%s) = %s, want %s", test.path, actual, test.expected)
			}
		})
	}
}

func TestDirPath_Matches_Empty(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	dirB := helpers.CreateTempDir(t)
	if m := dirA.Matches(dirB); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			dirA, dirB, paths.Match, m,
		)
	}
}

func TestDirPath_Matches_Itself(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)
	if m := dirA.Matches(dirA); m != paths.Match {
		t.Errorf(
			"Expected directories %s to match itself got: %s",
			dirA, m,
		)
	}
}

func TestDirPath_Matches_Unique(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)

	dirB := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirB)

	if m := dirA.Matches(dirB); m != paths.Mismatch {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			dirA, dirB, paths.Mismatch, m,
		)
	}
}

func TestDirPath_Copy(t *testing.T) {
	existingDir := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, existingDir)
	helpers.CreateSampleNestedStructure(t, paths.DirPath(filepath.Join(existingDir.FullPath(), "another")))

	tests := []struct {
		source         paths.DirPath
		dest           paths.DirPath
		copyShouldFail bool
		expects        paths.MatchType
		description    string
	}{
		{
			"not_found",
			"~/will_not_exist",
			true,
			paths.Mismatch,
			"Cannot copy non existent",
		},
		{
			helpers.CreateTempDir(t),
			helpers.CreateTempDirInHome(t),
			false,
			paths.Match,
			"Expands dest paths",
		},
		{
			helpers.CreateTempDirInHome(t),
			helpers.CreateTempDir(t),
			false,
			paths.Match,
			"Expands source paths",
		},
		{
			paths.DirPath(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			helpers.CreateTempDir(t),
			false,
			paths.Match,
			"Handles trailing slashes in src",
		},
		{
			helpers.CreateTempDir(t),
			paths.DirPath(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			false,
			paths.Match,
			"Handles trailing slashes in dest",
		},
		{
			paths.DirPath(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			paths.DirPath(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			false,
			paths.Match,
			"Handles trailing slashes in src and dest",
		},
		{
			helpers.CreateTempDir(t),
			existingDir,
			false,
			paths.Match,
			"Overwrites dest",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if exists, _ := test.source.Exists(); exists {
				helpers.CreateSampleNestedStructure(t, test.source)
			}

			if err := test.source.Copy(test.dest); err != nil && !test.copyShouldFail {
				t.Fatalf("Error copying a: %s to b: %s, err: %s", test.source, test.dest, err)
			}

			if m := test.source.Matches(test.dest); m != test.expects {
				t.Fatalf(
					"path mismatch a: %s, b: %s, want: %s, got: %s", test.source, test.dest, test.expects, m,
				)
			}
		})
	}
}
