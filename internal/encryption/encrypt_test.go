package encryption

import (
	"github.com/cshtdd/truecrypt/internal/paths"
	"github.com/cshtdd/truecrypt/internal/settings"
	"github.com/cshtdd/truecrypt/internal/test/helpers"
	"strings"
	"testing"
)

func TestEncrypt_SuccessOutput(t *testing.T) {
	lines := []struct {
		expected bool
		line     string
	}{
		{true, "Enter encryption password:"},
		{true, "Confirm password:"},
		{false, "secret1"},
	}
	for _, test := range lines {
		t.Run(test.line, func(t *testing.T) {
			f := helpers.NewFakeInput("secret1", "secret1")
			e := &encryptInput{
				in: f.In(),
				c:  newFakeZipper(),
				l: newFakeSettings(settings.Settings{
					DecryptedFolder: helpers.CreateTempDir(t),
					EncryptedFile:   helpers.CreateTempZip(t),
				}),
			}

			if err := encryptProgram(e); err != nil {
				t.Fatalf("Unexpected error. err: %s", err)
			}

			if found := strings.Contains(f.Out(), test.line); found != test.expected {
				t.Errorf("Line found got: %t want: %t output: %s", found, test.expected, f.Out())
			}
		})
	}
}

func TestEncrypt_FailsOnSettingsLoadError(t *testing.T) {
	s := newFakeSettingsWithError()
	e := &encryptInput{
		in: helpers.NewFakeInput("password", "mismatch").In(),
		l:  s,
	}

	if err := encryptProgram(e); err != s.Err() {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnInvalidEncryptedFile(t *testing.T) {
	e := &encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: helpers.CreateTempDir(t),
			EncryptedFile:   paths.ZipPath(helpers.CreateTemp(t)),
		}),
	}

	if err := encryptProgram(e); err == nil || err.Error() != "invalid encrypted file" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnInvalidDecryptedFolder(t *testing.T) {
	e := &encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: "not_found",
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}

	if err := encryptProgram(e); err == nil || err.Error() != "decrypted folder does not exist" {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEncrypt_FailsOnPasswordMismatch(t *testing.T) {
	d := helpers.CreateTempDir(t)
	e := &encryptInput{
		in: helpers.NewFakeInput("password", "mismatch").In(),
		l: newFakeSettings(settings.Settings{
			DecryptedFolder: d,
			EncryptedFile:   helpers.CreateTempZip(t),
		}),
	}

	if err := encryptProgram(e); err == nil || err.Error() != "passwords mismatch" {
		t.Fatalf("Unexpected error %s", err)
	}

	helpers.EnsureExists(t, d, true)
}

func TestEncrypt_FailsOnCompressionFailure(t *testing.T) {
	z := newFakeZipperWithError()
	s := settings.Settings{
		DecryptedFolder: helpers.CreateTempDir(t),
		EncryptedFile:   helpers.CreateTempZip(t),
	}
	e := &encryptInput{
		in: helpers.NewFakeInput("password", "password").In(),
		l:  newFakeSettings(s),
		c:  z,
	}

	if err := encryptProgram(e); err != z.Err() {
		t.Fatalf("Unexpected error %s", err)
	}

	helpers.EnsureExists(t, s.EncryptedFile, false)
	helpers.EnsureExists(t, s.DecryptedFolder, true)
}
