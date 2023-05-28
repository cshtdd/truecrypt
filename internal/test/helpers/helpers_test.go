package helpers_test

import (
	"fmt"
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

func TestPathComparisonEmptyDirectoriesMatch(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	dirB := helpers.CreateTempDir(t)
	if m := helpers.SampleDirectoriesMatch(t, dirA, dirB); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %d got: %d",
			dirA, dirB, paths.Match, m,
		)
	}
}

func TestPathComparisonAlwaysMatchItself(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)
	if m := helpers.SampleDirectoriesMatch(t, dirA, dirA); m != paths.Match {
		t.Errorf(
			"Expected directories %s to match itself got: %d",
			dirA, m,
		)
	}
}

func TestPathComparisonSamplesAreAlwaysUnique(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)

	dirB := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirB)

	if m := helpers.SampleDirectoriesMatch(t, dirA, dirB); m != paths.Mismatch {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %d got: %d",
			dirA, dirB, paths.Mismatch, m,
		)
	}
}

// These tests are here to avoid cyclical dependencies in the test
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
			helpers.CreateSampleNestedStructure(t, test.parent)
		})
	}
}

func TestPathDelete(t *testing.T) {
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

func TestPathMoveFile(t *testing.T) {
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
