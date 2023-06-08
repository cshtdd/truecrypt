package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
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

	if !s.EncryptedFile.IsValid() { // TODO: test the negative paths
		return errors.New("invalid encrypted file")
	}

	switch exists, err := s.DecryptedFolder.Exists(); { // TODO: test the negative paths
	case err != nil:
		return err
	case !exists:
		return errors.New("decrypted folder does not exist")
	}

	password, err := readPassword(e.in)
	if err != nil {
		return err
	}

	// ask for password confirmation
	e.in.WriteLine("Confirm password:")
	switch read, line, err := e.in.ReadLine(); {
	case err != nil:
		return err
	case !read || line != password:
		return errors.New("passwords mismatch")
	}
	// TODO: test password mismatch

	// TODO: test zip failure
	if err := e.c.compress(s, password); err != nil {
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
