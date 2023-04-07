package setup_test

import (
	"bytes"
	"strings"
	"testing"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestRunsSavesSettings(t *testing.T) {
	tests := []struct {
		settingsPath  paths.Path
		encryptedPath paths.Path
		description   string
	}{
		{helpers.CreateTemp(t), helpers.CreateTemp(t), "Encrypted file exists with full path"},
		{helpers.CreateTemp(t), helpers.CreateTempInHome(t), "Encrypted file exists on home"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var fakeOut bytes.Buffer
			input := setup.Input{
				IO: internal.IO{
					Reader: strings.NewReader(test.encryptedPath.String()),
					Writer: &fakeOut,
				},
				SettingsPath: test.settingsPath,
			}

			// run the program
			if err := setup.Run(input); err != nil {
				t.Error("Run failed", err)
			}

			// validate program prompts
			lines := strings.Join([]string{
				"Enter encrypted file:",
				test.encryptedPath.String(),
				"",
			}, "\n")
			if out := fakeOut.String(); out != lines {
				t.Errorf("Run() = %s, want %s", out, lines)
			}

			// validate saved settings
			savedSettings, err := settings.LoadFrom(test.settingsPath)
			if err != nil {
				t.Fatal("Cannot read saved settings", err)
			}

			if savedSettings.EncryptedFile != test.encryptedPath {
				t.Errorf(
					"savedSettings.EncryptedFile = %s, want %s",
					savedSettings.EncryptedFile, test.encryptedPath,
				)
			}
		})
	}
}

func TestRunFailsWhenEncryptedFileNotFound(t *testing.T) {
	tests := []struct {
		file        string
		description string
	}{
		{"not_found.zip", "File does not exist"},
		{"", "Blank input"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var fakeOut bytes.Buffer
			input := setup.Input{
				IO: internal.IO{
					Reader: strings.NewReader(test.file),
					Writer: &fakeOut,
				},
				SettingsPath: "does_not_matter.json",
			}
			// run the program
			if err := setup.Run(input); err == nil {
				t.Error("Expected Error")
			}
		})
	}
}
