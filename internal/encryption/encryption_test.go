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
	// create a test directory with many files
	decryptedSource := helpers.CreateTempDir(t)
	helpers.CreateSampleNestedStructure(t, decryptedSource)

	// create settings that point to that
	settingsEncrypt := helpers.CreateTemp(t)
	configEncrypt := settings.Settings{
		DecryptedFolder: decryptedSource,
		EncryptedFile:   helpers.CreateTemp(t),
	}
	if err := configEncrypt.Save(settingsEncrypt); err != nil {
		t.Fatalf("Error saving encrypt settings, err: %s", err)
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
		SettingsPath: settingsEncrypt,
	}
	if err := encryption.Encrypt(inputEncrypt); err != nil {
		t.Fatalf("Error encrypting data, err: %s", err)
	}

	// create a settings to decrypt to a different place
	settingsDecrypt := helpers.CreateTemp(t)
	configDecrypt := settings.Settings{
		DecryptedFolder: helpers.CreateTempDir(t),
		EncryptedFile:   configEncrypt.EncryptedFile,
	}
	if err := configDecrypt.Save(settingsDecrypt); err != nil {
		t.Fatalf("Error saving decrypt settings, err: %s", err)
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
		SettingsPath: settingsDecrypt,
	}
	if err := encryption.Decrypt(inputDecrypt); err != nil {
		t.Fatalf("Error decrypting data, err: %s", err)
	}

	// compare the two decrypted folders
	if m := helpers.SampleDirectoriesMatch(t, decryptedSource, configDecrypt.DecryptedFolder); m != paths.Match {
		t.Errorf(
			"Expected directories a: %s, b: %s to: %s got: %s",
			decryptedSource, configDecrypt.DecryptedFolder, paths.Match, m,
		)
	}
}
