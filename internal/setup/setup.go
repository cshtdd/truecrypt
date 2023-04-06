package setup

import (
	"bufio"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

type Input struct {
	internal.IO
	SettingsPath string
}

// func Run(io internal.IO) error {
func Run(in Input) error {
	// TODO: implement this
	in.IO.Writeln("Enter encrypted file:")

	var encryptedFile string

	// TODO: refactor this out to a nicer interface
	scanner := bufio.NewScanner(in.IO.Reader)
	if scanner.Scan() {
		encryptedFile = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	in.IO.Writeln(encryptedFile)

	s := settings.Settings{EncryptedFile: encryptedFile}
	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}

	return nil
}
