package setup

import (
	"bufio"
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

type Input struct {
	internal.IO
	SettingsPath paths.Path
}

// func Run(io internal.IO) error {
func Run(in Input) error {
	s := settings.Settings{}

	fmt.Fprintln(in.IO.Writer, "Enter encrypted file:")

	// TODO: refactor this out to a nicer interface
	scanner := bufio.NewScanner(in.IO.Reader)
	if scanner.Scan() {
		s.EncryptedFile = paths.Path(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if exists, err := s.EncryptedFile.Exists(); !exists {
		return err
	}

	// TODO: print the settings struct here
	fmt.Fprintln(in.IO.Writer, s.EncryptedFile)

	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}

	return nil
}
