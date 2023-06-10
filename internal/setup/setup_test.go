package setup_test

import (
	"fmt"
	"strings"
	"testing"

	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestRunsOutputValidation(t *testing.T) {
	tests := []struct {
		settingsPath  paths.FilePath
		encryptedPath paths.ZipPath
		description   string
	}{
		{
			helpers.CreateTemp(t),
			helpers.CreateTempZip(t),
			"Encrypted file exists with full path",
		},
		{
			helpers.CreateTemp(t),
			helpers.CreateTempZipInHome(t),
			"Encrypted file exists on home",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			f := helpers.NewFakeInputWithSettingsPath(
				test.settingsPath, test.encryptedPath.String(),
			)

			// run the program
			if err := setup.Run(f.In()); err != nil {
				t.Error("Run failed", err)
			}

			// validate program prompts
			expectedSettings := settings.Settings{
				DecryptedFolder: settings.DefaultDecryptedFolder(),
				EncryptedFile:   test.encryptedPath,
			}
			expectedLines := []string{
				"Enter encrypted file:\n",
				"Enter decrypted folder:\n",
				fmt.Sprintf("\tassumes \"%s\" if blank\n", settings.DefaultDecryptedFolder()),
				fmt.Sprintf("Saving settings: %+v\n", expectedSettings),
				fmt.Sprintf("Settings path: %s\n", test.settingsPath.String()),
			}
			for _, line := range expectedLines {
				if !strings.Contains(f.Out(), line) {
					t.Fatalf("Run() = %s, want %s", f.Out(), line)
				}
			}
		})
	}
}

func TestRunsSavesEncryptedPath(t *testing.T) {
	tests := []struct {
		settingsPath  paths.FilePath
		encryptedPath paths.ZipPath
		description   string
	}{
		{
			helpers.CreateTemp(t),
			helpers.CreateTempZip(t),
			"Encrypted file exists with full path",
		},
		{
			helpers.CreateTemp(t),
			helpers.CreateTempZipInHome(t),
			"Encrypted file exists on home",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			f := helpers.NewFakeInputWithSettingsPath(
				test.settingsPath, test.encryptedPath.String(),
			)

			// run the program
			if err := setup.Run(f.In()); err != nil {
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
		settingsPath            paths.FilePath
		inputDecryptedFolder    paths.DirPath
		expectedDecryptedFolder paths.DirPath
		description             string
	}{
		{
			helpers.CreateTemp(t),
			"",
			settings.DefaultDecryptedFolder(),
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
			f := helpers.NewFakeInputWithSettingsPath(
				test.settingsPath,
				helpers.CreateTempZip(t).String(), test.inputDecryptedFolder.String(),
			)

			// run the program
			if err := setup.Run(f.In()); err != nil {
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
		file        paths.ZipPath
		description string
	}{
		{"not_found.zip", "File does not exist"},
		{"", "Blank input"},
		{paths.ZipPath(helpers.CreateTemp(t)), "Invalid zip path"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			f := helpers.NewFakeInput(test.file.String())
			// run the program
			if err := setup.Run(f.In()); err == nil {
				t.Error("Expected Error")
			}
		})
	}
}

func TestLoadsExistingSettings(t *testing.T) {
	base := settings.Settings{
		EncryptedFile:   helpers.CreateTempZip(t),
		DecryptedFolder: paths.DirPath("~/tmp/decrypted"),
	}
	diffEncryptedFile := settings.Settings{
		EncryptedFile:   helpers.CreateTempZip(t),
		DecryptedFolder: paths.DirPath("~/tmp/decrypted"),
	}
	diffDecryptedFolder := settings.Settings{
		EncryptedFile:   base.EncryptedFile,
		DecryptedFolder: paths.DirPath("~/tmp/changed"),
	}

	tests := []struct {
		input       settings.Settings
		inputLines  []string
		expected    settings.Settings
		description string
	}{
		{base, []string{"", ""}, base, "Maintains settings untouched"},
		{base, []string{diffEncryptedFile.EncryptedFile.String(), ""}, diffEncryptedFile, "Overrides encrypted file"},
		{base, []string{"", diffDecryptedFolder.DecryptedFolder.String()}, diffDecryptedFolder, "Overrides decrypted folder"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// seed settings
			settingsPath := helpers.CreateTemp(t)
			if err := test.input.Save(settingsPath); err != nil {
				t.Fatalf("Unexpected error err: %s", err)
			}

			// fake user input
			f := helpers.NewFakeInputWithSettingsPath(settingsPath, test.inputLines...)

			// run the program
			if err := setup.Run(f.In()); err != nil {
				t.Error("Run failed", err)
			}

			// validate saved settings
			actual, err := settings.LoadFrom(settingsPath)
			if err != nil {
				t.Fatal("Cannot read saved settings", err)
			}

			if actual != test.expected {
				t.Errorf("Saved settings don't match. want = %+v, got = %+v", test.expected, actual)
			}

			// validate output
			expectedLines := []string{
				fmt.Sprintf("Loading settings from: %s\n", settingsPath),
				fmt.Sprintf("Enter encrypted file:\n\tassumes \"%s\" if blank\n", test.input.EncryptedFile),
				fmt.Sprintf("Enter decrypted folder:\n\tassumes \"%s\" if blank\n", test.input.DecryptedFolder),
			}
			for _, line := range expectedLines {
				if !strings.Contains(f.Out(), line) {
					t.Fatalf("Output contain mismatch. want = %s, got = %s", line, f.Out())
				}
			}
		})
	}
}

func TestRunFailsWhenDecryptedFolderIsBlank(t *testing.T) {
	// seed settings with no decrypted folder
	settingsPath := helpers.CreateTemp(t)
	s := settings.Settings{
		DecryptedFolder: "",
		EncryptedFile:   helpers.CreateTempZip(t),
	}
	if err := s.Save(settingsPath); err != nil {
		t.Fatalf("Error saving settings err: %s", err)
	}

	f := helpers.NewFakeInputWithSettingsPath(settingsPath, helpers.CreateTemp(t).String())

	if err := setup.Run(f.In()); err == nil {
		t.Error("Expected Error")
	}
}
