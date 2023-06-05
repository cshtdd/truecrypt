package settings_test

import (
	"os"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestDefaultSettingsPath(t *testing.T) {
	expected := paths.FilePath("~/.config/truecrypt/settings.json")
	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultSettingsPathEnvOverride(t *testing.T) {
	expected := paths.FilePath("aaa/bbb.json")
	os.Setenv("TC_SETTINGS", expected.String())
	defer os.Unsetenv("TC_SETTINGS")

	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolder(t *testing.T) {
	expected := paths.DirPath("~/decrypted_folder/t")
	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolderEnvOverride(t *testing.T) {
	expected := paths.DirPath("aaa/bbb")
	os.Setenv("TC_DECRYPTED", expected.String())
	defer os.Unsetenv("TC_DECRYPTED")

	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}

func TestSettingsSerialization(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	s := settings.Settings{DecryptedFolder: "aaa/bbb", EncryptedFile: "ccc.zip"}
	s.Save(tmp)

	x, err := settings.LoadFrom(tmp)
	if err != nil {
		t.Fatal("Got unexpected error", err)
	}

	if x != s {
		t.Errorf("got %v+, want %v+", x, s)
	}
}
