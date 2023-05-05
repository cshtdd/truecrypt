package settings_test

import (
	"os"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestDefaultSettingsPath(t *testing.T) {
	expected := "~/.config/truecrypt/settings.json"
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
	expected := "~/decrypted_folder/t"
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

func TestSettingsIsValidEncryptedFile(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	tests := []struct {
		file        paths.Path
		want        bool
		description string
	}{
		{"", false, "Empty file does not exist"},
		{"not_found", false, "Not found file does not exist"},
		{tmp, true, "File exists"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			s := settings.Settings{EncryptedFile: test.file}
			if got, _ := s.EncryptedFile.Exists(); got != test.want {
				t.Errorf("s.IsValidEncryptedFile(%s) = %t, want %t", test.file, got, test.want)
			}
		})
	}
}
