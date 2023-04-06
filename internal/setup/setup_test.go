package setup_test

import (
	"bytes"
	"strings"
	"testing"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestRunsSavesSettings(t *testing.T) {
	tmpSettings, tmpEncrypted := helpers.CreateTemp(t), helpers.CreateTemp(t)

	var fakeOut bytes.Buffer
	input := setup.Input{
		IO: internal.IO{
			Reader: strings.NewReader(tmpEncrypted),
			Writer: &fakeOut,
		},
		SettingsPath: tmpSettings,
	}

	// run the program
	if err := setup.Run(input); err != nil {
		t.Error("Run failed", err)
	}

	// validate program prompts
	lines := strings.Join([]string{
		"Enter encrypted file:",
		tmpEncrypted,
		"",
	}, "\n")
	if out := fakeOut.String(); out != lines {
		t.Errorf("Run() = %s, want %s", out, lines)
	}

	// validate saved settings
	savedSettings, err := settings.LoadFrom(tmpSettings)
	if err != nil {
		t.Fatal("Cannot read saved settings", err)
	}

	if savedSettings.EncryptedFile != tmpEncrypted {
		t.Errorf("savedSettings.EncryptedFile = %s, want %s", savedSettings.EncryptedFile, tmpEncrypted)
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
