package setup

import (
	"errors"
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

type Input struct {
	internal.IO
	SettingsPath paths.Path
}

func loadSettings(in Input, defaultSettings settings.Settings) (existing bool, s settings.Settings, err error) {
	exists, err := in.SettingsPath.Exists()
	if err != nil {
		return false, defaultSettings, err
	}

	if !exists {
		return false, defaultSettings, nil
	}

	existingSettings, err := settings.LoadFrom(in.SettingsPath)
	if err != nil {
		return false, defaultSettings, err
	}

	var emptySettings settings.Settings
	if existingSettings == emptySettings {
		return false, defaultSettings, nil
	}

	return true, existingSettings, nil
}

func Run(in Input) error {
	// loading settings
	defaultSettings := settings.Settings{
		DecryptedFolder: paths.Path(settings.DefaultDecryptedFolder()),
	}
	existingSettings, s, err := loadSettings(in, defaultSettings)
	if err != nil {
		return err
	}
	if existingSettings {
		in.WriteLine(fmt.Sprintf("Loading settings from: %s", in.SettingsPath))
	}

	// reading encrypted file
	in.WriteLine("Enter encrypted file:")
	if existingSettings {
		in.WriteLine(fmt.Sprintf("\tassumes \"%s\" if blank", s.EncryptedFile))
	}
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read:
		s.EncryptedFile = paths.Path(line)
	}
	if exists, err := s.EncryptedFile.Exists(); !exists {
		return err
	}

	// reading decrypted folder
	defaultDecryptedFolder := defaultSettings.DecryptedFolder
	if existingSettings {
		defaultDecryptedFolder = s.DecryptedFolder
	}
	in.WriteLine("Enter decrypted folder:")
	in.WriteLine(fmt.Sprintf("\tassumes \"%s\" if blank", defaultDecryptedFolder))
	switch read, line, err := in.ReadLine(); {
	case err != nil:
		return err
	case read:
		s.DecryptedFolder = paths.Path(line)
	}
	if len(s.DecryptedFolder) == 0 {
		return errors.New("decrypted folder cannot be blank")
	}

	// savings settings
	in.WriteLine(fmt.Sprintf("Saving settings: %+v", s))
	in.WriteLine(fmt.Sprintf("Settings Path: %s", in.SettingsPath.String()))
	if err := s.Save(in.SettingsPath); err != nil {
		return err
	}
	return nil
}
