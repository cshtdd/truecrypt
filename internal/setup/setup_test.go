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

func TestRunsOutputValidation(t *testing.T) {
	tests := []struct {
		settingsPath    paths.Path
		encryptedPath   paths.Path
		decryptedFolder paths.Path
		description     string
	}{
		{
			helpers.CreateTemp(t),
			helpers.CreateTemp(t),
			"",
			"Encrypted file exists with full path",
		},
		{
			helpers.CreateTemp(t),
			helpers.CreateTempInHome(t),
			"",
			"Encrypted file exists on home",
		},
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
				"Enter decrypted folder:",
				test.encryptedPath.String(),
				"",
			}, "\n")
			if out := fakeOut.String(); out != lines {
				t.Errorf("Run() = %s, want %s", out, lines)
			}
		})
	}
}

func TestRunsSavesEncryptedPath(t *testing.T) {
	tests := []struct {
		settingsPath  paths.Path
		encryptedPath paths.Path
		description   string
	}{
		{
			helpers.CreateTemp(t),
			helpers.CreateTemp(t),
			"Encrypted file exists with full path",
		},
		{
			helpers.CreateTemp(t),
			helpers.CreateTempInHome(t),
			"Encrypted file exists on home",
		},
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

func TestRunsSavesDecryptedFolder(t *testing.T) {
	tests := []struct {
		settingsPath            paths.Path
		inputDecryptedFolder    string
		expectedDecryptedFolder paths.Path
		description             string
	}{
		{
			helpers.CreateTemp(t),
			"",
			paths.Path(settings.DefaultDecryptedFolder()),
			"Assumes default",
		},
		{
			helpers.CreateTemp(t),
			"~/tmp/decrypted",
			"~/tmp/decrypted",
			"Saves input",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var fakeOut bytes.Buffer
			input := setup.Input{
				IO: internal.IO{
					Reader: strings.NewReader(
						strings.Join([]string{helpers.CreateTemp(t).String(), test.inputDecryptedFolder}, "\n"),
					),
					Writer: &fakeOut,
				},
				SettingsPath: test.settingsPath,
			}

			// run the program
			if err := setup.Run(input); err != nil {
				t.Error("Run failed", err)
			}

			// validate saved settings
			savedSettings, err := settings.LoadFrom(test.settingsPath)
			if err != nil {
				t.Fatal("Cannot read saved settings", err)
			}

			if savedSettings.DecryptedFolder != test.expectedDecryptedFolder {
				t.Errorf(
					"savedSettings.DecryptedFolder = %s, want %s",
					savedSettings.DecryptedFolder,
					test.expectedDecryptedFolder,
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
