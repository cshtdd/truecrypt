package encryption

import (
	"errors"
	"fmt"
	"os/exec"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

func Encrypt(in internal.Input) error {
	in.WriteLine(fmt.Sprintf("Loading settings from: %s", in.SettingsPath))
	s, err := settings.LoadFrom(in.SettingsPath)
	if err != nil {
		return err
	}
	in.WriteLine(fmt.Sprintf("Loaded settings: %+v", s))

	// validate decrypted folder is valid
	exists, err := s.DecryptedFolder.Exists() // TODO: test the negative paths
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("decrypted folder does not exist")
	}

	// ask for password
	var password string
	in.WriteLine("Enter encryption password:")
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read && len(line) > 5:
		password = line
	default:
		return errors.New("empty or short passwords are not allowed")
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
	return compress(s.DecryptedFolder, s.EncryptedFile, password)
}

func Decrypt(in internal.Input) error {
	// TODO implement this
	return nil
}

func compress(sourceDir paths.Path, dest paths.Path, password string) error {
	// TODO: maybe figure out a way to inject this
	// TODO: ensure zip can run
	tmp, err := paths.CreateTempFile()
	if err != nil {
		return err
	}

	cmd := exec.Command("zip", "-e", "-r", "-P", password, tmp.FullPath(), sourceDir.FullPath())
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return tmp.MoveFile(dest)
}
