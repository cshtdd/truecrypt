package setup

import (
	"errors"
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

func loadSettings(in internal.Input, defaultSettings settings.Settings) (existing bool, s settings.Settings) {
	exists, err := in.SettingsPath.Exists()
	if !exists || err != nil {
		return false, defaultSettings
	}

	var emptySettings settings.Settings
	existingSettings, err := settings.LoadFrom(in.SettingsPath)
	if err != nil || existingSettings == emptySettings {
		return false, defaultSettings
	}

	return true, existingSettings
}

func Run(in internal.Input) error {
	// loading settings
	defaultSettings := settings.Settings{
		DecryptedFolder: paths.Path(settings.DefaultDecryptedFolder()),
	}
	existingSettings, s := loadSettings(in, defaultSettings)
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
