package encryption

import (
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
	"testing"
)

func TestDecrypt_FailsOnSettingsLoadError(t *testing.T) {
	s := newFakeSettingsWithError()
	d := decryptInput{
		in: helpers.NewFakeInput("password").In(),
		l:  s,
	}
	if err := decryptProgram(&d); err != s.Err() {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestDecrypt_FailsOnNotFoundEncryptedFile(t *testing.T) {
	d := decryptInput{
		in: helpers.NewFakeInput("password").In(),
		e:  nil,
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   "not_found.zip",
		}),
	}
	if err := decryptProgram(&d); err == nil || err.Error() != "encrypted file does not exist" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestDecrypt_FailsOnInvalidEncryptedFile(t *testing.T) {
	d := decryptInput{
		in: helpers.NewFakeInput("password").In(),
		e:  nil,
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   paths.ZipPath(helpers.CreateTemp(t)),
		}),
	}
	if err := decryptProgram(&d); err == nil || err.Error() != "invalid encrypted file" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestDecrypt_FailsOnExistingDecryptedFolder(t *testing.T) {
	d := decryptInput{
		in: helpers.NewFakeInput("password").In(),
		e:  nil,
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}
	if err := decryptProgram(&d); err == nil || err.Error() != "decrypted folder exists" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestDecrypt_FailsOnExtractionFailure(t *testing.T) {
	f := helpers.CreateTempDir(t)
	if err := f.Delete(); err != nil {
		t.Fatalf("Unexpected error deleting folder. err: %s", err)
	}

	z := newFakeZipperWithError()
	d := decryptInput{
		in: helpers.NewFakeInput("password").In(),
		e:  z,
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: f,
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}
	if err := decryptProgram(&d); err != z.Err() {
		t.Fatalf("Unexpected error %s", err)
	}

	helpers.EnsureExists(t, f, false)
}
