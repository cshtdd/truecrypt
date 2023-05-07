package clean

import (
	"fmt"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

func Run(in internal.Input) error {
	existingSettings, err := settings.LoadFrom(in.SettingsPath)
	if err != nil {
		return err
	}
	in.WriteLine(fmt.Sprintf("Deleting decrypted folder: %s", existingSettings.DecryptedFolder))
	return existingSettings.DecryptedFolder.Delete()
}

func CleanSettings(in internal.Input) error {
	in.WriteLine(fmt.Sprintf("Deleting settings at: %s", in.SettingsPath))
	return in.SettingsPath.Delete()
}
