package encryption

import (
	"errors"
	"github.com/cshtdd/truecrypt/internal"
	"github.com/cshtdd/truecrypt/internal/settings"
)

// Encrypt is the encryption program
func Encrypt(in *internal.Input) error {
	p := &encryptInput{
		in: in,
		c:  newZipper(in),
	}
	p.l = p
	return encryptProgram(p)
}

// private implementation with test shims
func encryptProgram(e *encryptInput) error {
	s, err := e.l.load()
	if err != nil {
		return err
	}

	if !s.EncryptedFile.IsValid() {
		return errors.New("invalid encrypted file")
	}

	switch exists, err := s.DecryptedFolder.Exists(); {
	case err != nil:
		return err
	case !exists:
		return errors.New("decrypted folder does not exist")
	}

	password, err := readPassword(e.in)
	if err != nil {
		return err
	}

	if err := confirmPassword(e.in, password); err != nil {
		return err
	}

	if err := e.c.compress(s, password); err != nil {
		_ = s.EncryptedFile.Delete() // delete the encrypted file to avoid trash
		return err
	}

	// wipe the original directory on success
	return s.DecryptedFolder.Delete()
}

// test shims

type encryptInput struct {
	in *internal.Input
	c  compressor
	l  settingsLoader
}

func (e *encryptInput) load() (settings.Settings, error) {
	return loadSettings(e.in)
}
