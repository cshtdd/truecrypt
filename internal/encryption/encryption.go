package encryption

import (
	"errors"
	"fmt"
	"os/exec"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

func loadSettings(in internal.Input) (settings.Settings, error) {
	in.WriteLine(fmt.Sprintf("Loading settings from: %s", in.SettingsPath))
	s, err := settings.LoadFrom(in.SettingsPath)
	if err == nil {
		in.WriteLine(fmt.Sprintf("Loaded settings: %+v", s))
	}
	return s, err
}

func readPassword(in *internal.Input) (string, error) {
	in.WriteLine("Enter encryption password:")
	// TODO: make sure ReadLine is not such a diva that needs a pointer, or pass an in pointer everywhere
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return "", err
	case read && len(line) > 5:
		return line, nil
	default:
		return "", errors.New("empty or short passwords are not allowed")
	}
}

func Encrypt(in internal.Input) error {
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

	password, err := readPassword(&in)
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
	if err := compress(s.DecryptedFolder, s.EncryptedFile, password); err != nil {
		return err
	}

	// wipe the original directory on success
	return s.DecryptedFolder.Delete()
}

func Decrypt(in internal.Input) error {
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
	return extract(s.EncryptedFile, s.DecryptedFolder, password)
}

func compress(sourceDir paths.DirPath, dest paths.ZipPath, password string) error {
	// TODO: maybe figure out a way to inject this method
	if _, err := exec.Command("zip", "-v").Output(); err != nil { // TODO display this output
		return err
	}

	tmp, err := paths.CreateTempZipFile()
	if err != nil {
		return err
	}

	cmd := exec.Command("zip", "-er", "-P", password, tmp.FullPath(), sourceDir.Base())
	cmd.Dir = sourceDir.DirName()
	if _, err := cmd.Output(); err != nil { // TODO display this output
		return err
	}

	return tmp.Move(dest)
}

func extract(source paths.ZipPath, dest paths.DirPath, password string) error {
	// TODO: maybe figure out a way to inject this method
	if _, err := exec.Command("unzip", "-v").Output(); err != nil { // TODO display this output
		return err
	}

	if err := dest.Create(); err != nil {
		return err
	}

	cmd := exec.Command("unzip", "-n", source.FullPath(), "-d", dest.FullPath(), "-P", password)
	_, err := cmd.Output()
	return err
}
