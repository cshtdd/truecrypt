package encryption_test

import (
	"path/filepath"
	"tddapps.com/truecrypt/internal/encryption"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
	"testing"
)

func TestEncryption_EndToEnd(t *testing.T) {
	tests := []struct {
		password    string
		description string
	}{
		{"simplePassword1", "Regular Password"},
		{"passwords with spaces", "Passphrases"},
		{"@,<>?/[]{}#$%^&*()+=", "All special chars"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// create settings that point to a decrypted folder
			sp := helpers.CreateTemp(t)
			zipInNestedDir := paths.ZipPath(filepath.Join(helpers.CreateTempDir(t).FullPath(), "folder1/folder2/a.zip"))
			s := settings.Settings{
				DecryptedFolder: helpers.CreateTempDir(t),
				EncryptedFile:   zipInNestedDir,
			}
			// create a test directory with many files
			helpers.CreateSampleNestedStructure(t, s.DecryptedFolder)
			// save the settings
			if err := s.Save(sp); err != nil {
				t.Fatalf("Error saving encryptProgram settings, err: %s", err)
			}

			// clone the source because it will get wiped
			clone := helpers.CreateTempDir(t)
			if err := s.DecryptedFolder.Copy(clone); err != nil {
				t.Fatalf("Unexpected error cloning source, err: %s", err)
			}
			if m := s.DecryptedFolder.Matches(clone); m != paths.Match {
				t.Fatalf("Expected clone to match, got: %s", m)
			}

			// encryptProgram the data
			fe := helpers.NewFakeInputWithSettingsPath(sp, test.password, test.password)
			if err := encryption.Encrypt(fe.In()); err != nil {
				t.Fatalf("Error encrypting data, err: %s", err)
			}

			// source must get wiped
			if exists, err := s.DecryptedFolder.Exists(); exists || err != nil {
				t.Fatalf("Decrypted source should not exist")
			}

			// decryptProgram the data
			fd := helpers.NewFakeInputWithSettingsPath(sp, test.password)
			if err := encryption.Decrypt(fd.In()); err != nil {
				t.Fatalf("Error decrypting data, err: %s", err)
			}

			// compare the two decrypted folders
			if m := clone.Matches(s.DecryptedFolder); m != paths.Match {
				t.Errorf(
					"Expected directories a: %s, b: %s to: %s got: %s",
					clone, s.DecryptedFolder, paths.Match, m,
				)
			}
		})
	}
}
