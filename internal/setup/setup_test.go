package setup_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
)

func TestRun(t *testing.T) {
	// setup a temp file for the saved settings
	tmpSettings, err := os.CreateTemp("", "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
	}
	defer os.Remove(tmpSettings.Name())

	// setup a temp file for the encrypted file
	tmpEncrypted, err := os.CreateTemp("", "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
	}
	defer os.Remove(tmpEncrypted.Name())

	// setup fake user inputs
	var fakeOut bytes.Buffer
	input := setup.Input{
		IO: internal.IO{
			Reader: strings.NewReader(tmpEncrypted.Name()),
			Writer: &fakeOut,
		},
		SettingsPath: tmpSettings.Name(),
	}
	// run the program
	if err := setup.Run(input); err != nil {
		t.Error("Run failed", err)
	}

	// validate program prompts
	lines := strings.Join([]string{
		"Enter encrypted file:",
		tmpEncrypted.Name(),
		"",
	}, "\n")
	if out := fakeOut.String(); out != lines {
		t.Errorf("Run() = %s, want %s", out, lines)
	}

	// validate saved settings
	savedSettings, err := settings.LoadFrom(tmpSettings.Name())
	if err != nil {
		t.Fatal("Cannot read saved settings", err)
	}

	if savedSettings.EncryptedFile != tmpEncrypted.Name() {
		t.Errorf("savedSettings.EncryptedFile = %s, want %s", savedSettings.EncryptedFile, tmpEncrypted.Name())
	}
}
