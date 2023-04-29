package setup

import (
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

type Input struct {
	internal.IO
	SettingsPath paths.Path
}

func Run(in Input) error {
	s := settings.Settings{
		DecryptedFolder: paths.Path(settings.DefaultDecryptedFolder()),
	}

	fmt.Fprintln(in.IO.Writer, "Enter encrypted file:")
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read:
		s.EncryptedFile = paths.Path(line)
	}
	if exists, err := s.EncryptedFile.Exists(); !exists {
		return err
	}

	fmt.Fprintln(in.IO.Writer, "Enter decrypted folder:")
	switch read, line, err := in.ReadLine(); {
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
