package encryption

import (
	"errors"
	"fmt"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

// Shared IO functions

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
