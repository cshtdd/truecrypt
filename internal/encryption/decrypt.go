package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

// Decrypt is the decrypt program
func Decrypt(in *internal.Input) error {
	p := &decryptInput{
		in: in,
		e:  newZipper(in),
	}
	p.l = p
	return decryptProgram(p)
}

// private implementation with test shims
func decryptProgram(d *decryptInput) error {
	s, err := d.l.load()
	if err != nil {
		return err
	}

	switch exists, err := s.EncryptedFile.Exists(); {
	case err != nil:
		return err
	case !exists:
		return errors.New("encrypted file does not exist")
	}
	if !s.EncryptedFile.IsValid() {
		return errors.New("invalid encrypted file")
	}

	switch exists, err := s.DecryptedFolder.Exists(); {
	case err != nil:
		return err
	case exists:
		return errors.New("decrypted folder exists")
	}

	password, err := readPassword(d.in)
	if err != nil {
		return err
	}

	return d.e.extract(s, password)
}

// test shims
type decryptInput struct {
	in *internal.Input
	e  extractor
	l  settingsLoader
}

func (d *decryptInput) load() (settings.Settings, error) {
	return loadSettings(d.in)
}
