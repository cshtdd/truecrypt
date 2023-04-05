package settings_test

import (
	"fmt"
	"os"
	"testing"

	"tddapps.com/truecrypt/internal/settings"
)

func TestDefaultSettingsPath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Cannot read pwd", err)
	}

	expected := fmt.Sprintf("%s/.settings.json", pwd)
	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultSettingsPathEnvOverride(t *testing.T) {
	expected := "aaa/bbb.json"
	os.Setenv("TC_SETTINGS", expected)
	defer os.Unsetenv("TC_SETTINGS")

	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolder(t *testing.T) {
	expected := fmt.Sprintf("%s/decrypted_folder/t", os.Getenv("HOME"))
	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolderEnvOverride(t *testing.T) {
	expected := "aaa/bbb.json"
	os.Setenv("TC_DECRYPTED", expected)
	defer os.Unsetenv("TC_DECRYPTED")

	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}

func TestSettingsSerialization(t *testing.T) {
	tmp, err := os.CreateTemp("", "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
	}
	defer os.Remove(tmp.Name())

	s := settings.Settings{DecryptedFolder: "aaa/bbb", EncryptedFile: "ccc.zip"}
	s.Save(tmp.Name())

	x, err := settings.LoadFrom(tmp.Name())
	if err != nil {
		t.Fatal("Got unexpected error", err)
	}

	if x != s {
		t.Errorf("got %v+, want %v+", x, s)
	}
}
