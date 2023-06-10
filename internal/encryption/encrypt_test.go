package encryption

import (
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
	"testing"
)

func TestEncrypt_FailsOnSettingsLoadError(t *testing.T) {
	s := newFakeSettingsWithError()
	e := encryptInput{
		in: helpers.NewFakeInput("password", "mismatch").In(),
		l:  s,
	}

	if err := encryptProgram(&e); err != s.Err() {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnInvalidEncryptedFile(t *testing.T) {
	e := encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   paths.ZipPath(helpers.CreateTemp(t)),
		}),
	}

	if err := encryptProgram(&e); err == nil || err.Error() != "invalid encrypted file" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnInvalidDecryptedFolder(t *testing.T) {
	e := encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: "not_found",
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}

	if err := encryptProgram(&e); err == nil || err.Error() != "decrypted folder does not exist" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnPasswordMismatch(t *testing.T) {
	e := encryptInput{
		in: helpers.NewFakeInput("password", "mismatch").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}

	if err := encryptProgram(&e); err == nil || err.Error() != "passwords mismatch" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnCompressionFailure(t *testing.T) {
	z := newFakeZipperWithError()
	e := encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
		c: z,
	}

	if err := encryptProgram(&e); err != z.Err() {
		t.Fatalf("Unexpected error %s", err)
	}
}
