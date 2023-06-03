package encryption_test

import (
	"bytes"
	"strings"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/encryption"
	"tddapps.com/truecrypt/internal/paths"
	"testing"

	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestEndToEnd(t *testing.T) {
	// create settings that point to a decrypted folder
	sp := helpers.CreateTemp(t)
	s := settings.Settings{
		DecryptedFolder: helpers.CreateTempDir(t),
		EncryptedFile:   helpers.CreateTemp(t),
	}
	// create a test directory with many files
	helpers.CreateSampleNestedStructure(t, s.DecryptedFolder)
	// save the settings
	if err := s.Save(sp); err != nil {
		t.Fatalf("Error saving encrypt settings, err: %s", err)
	}

	// clone the source because it will get wiped
	clone := helpers.CreateTempDir(t)
	if err := s.DecryptedFolder.Copy(clone); err != nil {
		t.Fatalf("Unexpected error cloning source, err: %s", err)
	}
	if m := s.DecryptedFolder.Matches(clone); m != paths.Match {
		t.Fatalf("Expected clone to match, got: %s", m)
	}

	password := "thisissecret"

	// encrypt the data
	var fakeOutEncrypt bytes.Buffer
	inputEncrypt := internal.Input{
		IO: internal.IO{
			Reader: strings.NewReader(
				strings.Join([]string{password, password, ""}, "\n"),
			),
			Writer: &fakeOutEncrypt,
		},
		SettingsPath: sp,
	}
	if err := encryption.Encrypt(inputEncrypt); err != nil {
		t.Fatalf("Error encrypting data, err: %s", err)
	}

	// source must get wiped
	if exists, err := s.DecryptedFolder.Exists(); exists || err != nil {
		t.Fatalf("Decrypted source should not exist")
	}

	// decrypt the data
	var fakeOutDecrypt bytes.Buffer
	inputDecrypt := internal.Input{
		IO: internal.IO{
			Reader: strings.NewReader(
				strings.Join([]string{password, ""}, "\n"),
			),
			Writer: &fakeOutDecrypt,
		},
		SettingsPath: sp,
	}
	if err := encryption.Decrypt(inputDecrypt); err != nil {
		t.Fatalf("Error decrypting data, err: %s", err)
	}

	// compare the two decrypted folders
	if m := clone.Matches(s.DecryptedFolder); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			clone, s.DecryptedFolder, paths.Match, m,
		)
	}
}
