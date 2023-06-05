package helpers_test

import (
	"path/filepath"
	"testing"

	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestCreateTemp(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	if exists, err := tmp.Exists(); !exists {
		t.Errorf("Tmp file does not exist %s, err %s", tmp, err)
	}
}

func TestCreateTempZip(t *testing.T) {
	tmp := helpers.CreateTempZip(t)

	if exists, err := tmp.Exists(); !exists {
		t.Errorf("Tmp Zip file does not exist %s, err %s", tmp, err)
	}

	if v := tmp.IsValid(); !v {
		t.Errorf("Tmp Zip file is not valid %s", tmp)
	}
}

func TestCreateTempInHome(t *testing.T) {
	tmp := helpers.CreateTempInHome(t)

	expected := filepath.Join("~", tmp.Base())
	if tmp.String() != expected {
		t.Errorf("CreateTempInHome() got %s, want %s", tmp, expected)
	}
}

func TestCreateTempZipInHome(t *testing.T) {
	tmp := helpers.CreateTempZipInHome(t)

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
