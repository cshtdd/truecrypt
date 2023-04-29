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

type lineScanner bufio.Scanner

func (s *lineScanner) readLine() (read bool, line string, err error) {
	scanner := (*bufio.Scanner)(s)
	if read = scanner.Scan(); read {
		line = scanner.Text()
	}
	return read, line, scanner.Err()
}

func Run(in Input) error {
	s := settings.Settings{
		DecryptedFolder: paths.Path(settings.DefaultDecryptedFolder()),
	}

	fmt.Fprintln(in.IO.Writer, "Enter encrypted file:")

	scanner := (*lineScanner)(bufio.NewScanner(in.IO.Reader))
	switch read, line, err := scanner.readLine(); {
	case err != nil:
		return err
	case read:
		s.EncryptedFile = paths.Path(line)
	}
	if exists, err := s.EncryptedFile.Exists(); !exists {
		return err
	}

	fmt.Fprintln(in.IO.Writer, "Enter decrypted folder:")
	switch read, line, err := scanner.readLine(); {
	case err != nil:
		return err
	case read:
		s.DecryptedFolder = paths.Path(line)
	}

	fmt.Fprintf(in.IO.Writer, "Saving settings: %+v\n", s)
	fmt.Fprintf(in.IO.Writer, "Settings Path: %s\n", in.SettingsPath.String())

	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}

	return nil
}
