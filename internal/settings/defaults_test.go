package settings_test

import (
	"github.com/cshtdd/truecrypt/internal/paths"
	"github.com/cshtdd/truecrypt/internal/settings"
	"github.com/cshtdd/truecrypt/internal/test/helpers"
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
	helpers.SetEnv("TC_SETTINGS", expected.String(), t)

	if p := settings.DefaultSettingsPath(); p != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolder(t *testing.T) {
	expected := paths.DirPath("~/decrypted_folder")
	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}

func TestDefaultDecryptedFolder_EnvOverride(t *testing.T) {
	expected := paths.DirPath("aaa/bbb")
	helpers.SetEnv("TC_DECRYPTED", expected.String(), t)

	if p := settings.DefaultDecryptedFolder(); p != expected {
		t.Errorf("DefaultDecryptedFolder() = %s, want %s", p, expected)
	}
}
