package settings_test

import (
	"os"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"testing"
)

func TestDefaultSettingsPath(t *testing.T) {
	expected := paths.FilePath("~/.config/truecrypt/settings.json")
	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultSettingsPath_EnvOverride(t *testing.T) {
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

func TestDefaultDecryptedFolder_EnvOverride(t *testing.T) {
	expected := paths.DirPath("aaa/bbb")
	os.Setenv("TC_DECRYPTED", expected.String())
	defer os.Unsetenv("TC_DECRYPTED")

	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}
