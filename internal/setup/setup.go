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

	in.WriteLine("Enter encrypted file:")
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read:
		s.EncryptedFile = paths.Path(line)
	}
	if exists, err := s.EncryptedFile.Exists(); !exists {
		return err
	}

	in.WriteLine("Enter decrypted folder:")
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read:
		s.DecryptedFolder = paths.Path(line)
	}

	in.WriteLine(fmt.Sprintf("Saving settings: %+v", s))
	in.WriteLine(fmt.Sprintf("Settings Path: %s", in.SettingsPath.String()))

	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}

	return nil
}
