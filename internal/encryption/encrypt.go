package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

type Compressor interface {
	Compress(s settings.Settings, password string) error
}

func Encrypt(in *internal.Input) error {
	return encrypt(in, newZipper(in))
}

func encrypt(in *internal.Input, c Compressor) error {
	s, err := loadSettings(in)
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

	password, err := readPassword(in)
	if err != nil {
		return err
	}

	// ask for password confirmation
	in.WriteLine("Confirm password:")
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case !read || line != password:
		return errors.New("passwords mismatch")
	}
	// TODO: test password mismatch

	// TODO: test zip failure
	if err := c.Compress(s, password); err != nil {
		return err
	}

	// wipe the original directory on success
	return s.DecryptedFolder.Delete()
}
