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

func TestPathComparisonEmptyDirectoriesMatch(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	dirB := helpers.CreateTempDir(t)
	if m := helpers.SampleDirectoriesMatch(t, dirA, dirB); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			dirA, dirB, paths.Match, m,
		)
	}
}

func TestPathComparisonAlwaysMatchItself(t *testing.T) {
	dirA := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, dirA)
	if m := helpers.SampleDirectoriesMatch(t, dirA, dirA); m != paths.Match {
		t.Errorf(
			"Expected directories %s to match itself got: %s",
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
