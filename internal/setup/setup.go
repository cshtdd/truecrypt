package setup

import (
	"bufio"
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

type Input struct {
	internal.IO
	SettingsPath string
}

// func Run(io internal.IO) error {
func Run(in Input) error {
	s := settings.Settings{}

	fmt.Fprintln(in.IO.Writer, "Enter encrypted file:")

	// TODO: refactor this out to a nicer interface
	scanner := bufio.NewScanner(in.IO.Reader)
	if scanner.Scan() {
		s.EncryptedFile = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if valid, err := s.IsValidEncryptedFile(); !valid {
		return err
	}

	fmt.Fprintln(in.IO.Writer, s.EncryptedFile)

	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}

	return nil
}
