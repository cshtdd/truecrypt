package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal"
)

// decrypt program

func decrypt(in internal.Input, e Extractor) error {
	s, err := loadSettings(in)
	if err != nil {
		return err
	}

	switch exists, err := s.EncryptedFile.Exists(); { // TODO: test the negative paths
	case err != nil:
		return err
	case !exists:
		return errors.New("encrypted file does not exist")
	}
	if !s.EncryptedFile.IsValid() { // TODO: test the negative paths
		return errors.New("invalid encrypted file")
	}

	switch exists, err := s.DecryptedFolder.Exists(); { // TODO: test the negative paths
	case err != nil:
		return err
	case exists:
		return errors.New("decrypted folder exists")
	}

	password, err := readPassword(&in)
	if err != nil {
		return err
	}

	// TODO: test unzip failure
	return e.Extract(s, password)
}
