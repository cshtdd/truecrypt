package paths_test

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestCreateDir(t *testing.T) {
	t1 := helpers.CreateTempDir(t)
	t2 := helpers.CreateTempDirInHome(t)

	tests := []struct {
		p           paths.Path
		target      paths.Path
		description string
	}{
		{
			paths.Path(fmt.Sprintf("%s/aaa.txt", t1)),
			t1,
			"Does not create existing dir",
		},
		{
			paths.Path(fmt.Sprintf("%s/bbb/ccc/aaa.txt", t1)),
			paths.Path(fmt.Sprintf("%s/bbb/ccc/", t1)),
			"Creates dir structure"},
		{
			paths.Path(fmt.Sprintf("%s/bbb/ccc/ddd/eee/", t1)),
			paths.Path(fmt.Sprintf("%s/bbb/ccc/ddd/eee/", t1)),
			"Creates dir path"},
		{
			paths.Path(fmt.Sprintf("%s/aaa.txt", t2)),
			t2, "Does not create existing home dir"},
		{
			paths.Path(fmt.Sprintf("%s/bbb/ccc/aaa.txt", t1)),
			paths.Path(fmt.Sprintf("%s/bbb/ccc/", t1)),
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

func TestMatchesDirEmptyDirectoriesMatch(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	dirB := helpers.CreateTempDir(t)
	if m := dirA.MatchesDir(dirB); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			dirA, dirB, paths.Match, m,
		)
	}
}

func TestMatchesDirAlwaysMatchItself(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)
	if m := dirA.MatchesDir(dirA); m != paths.Match {
		t.Errorf(
			"Expected directories %s to match itself got: %s",
			dirA, m,
		)
	}
}

func TestMatchesDirSamplesAreAlwaysUnique(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)

	dirB := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirB)

	if m := dirA.MatchesDir(dirB); m != paths.Mismatch {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			dirA, dirB, paths.Mismatch, m,
		)
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
			switch m, err := control.MatchesFile(test.dest); {
			case err != nil:
				t.Fatalf("Unexpected error comparing files a: %s, b:%s, err: %s", control, test.src, err)
			case m != paths.Match:
				t.Fatalf("Dest path doesn't match control")
			}
		})
	}
}

func TestCopyDirectory(t *testing.T) {
	existingDir := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, existingDir)
	helpers.CreateSampleNestedStructure(t, paths.Path(filepath.Join(existingDir.FullPath(), "another")))

	tests := []struct {
		source         paths.Path
		dest           paths.Path
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
			paths.Path(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			helpers.CreateTempDir(t),
			false,
			paths.Match,
			"Handles trailing slashes in src",
		},
		{
			helpers.CreateTempDir(t),
			paths.Path(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			false,
			paths.Match,
			"Handles trailing slashes in dest",
		},
		{
			paths.Path(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
			paths.Path(fmt.Sprintf("%s/", helpers.CreateTempDir(t))),
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

			if err := test.source.CopyDir(test.dest); err != nil && !test.copyShouldFail {
				t.Fatalf("Error copying a: %s to b: %s, err: %s", test.source, test.dest, err)
			}

			if m := test.source.MatchesDir(test.dest); m != test.expects {
				t.Fatalf(
					"Path mismatch a: %s, b: %s, want: %s, got: %s", test.source, test.dest, test.expects, m,
				)
			}
		})
	}
}
