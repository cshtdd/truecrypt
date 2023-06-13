package encryption

import (
	"errors"
	"fmt"
	"github.com/cshtdd/truecrypt/internal"
	"github.com/cshtdd/truecrypt/internal/settings"
)

// Shared IO functions

type settingsLoader interface {
	load() (settings.Settings, error)
}

func loadSettings(in *internal.Input) (settings.Settings, error) {
	in.WriteLine(fmt.Sprintf("Loading settings from: %s", in.SettingsPath))
	s, err := settings.LoadFrom(in.SettingsPath)
	if err == nil {
		in.WriteLine(fmt.Sprintf("Loaded settings: %+v", s))
	}
	return s, err
}

func readPassword(in *internal.Input) (string, error) {
	in.WriteLine("Enter encryption password:")
	switch read, line, err := in.ReadSensitiveLine(); {
	case err != nil:
		return "", err
	case read && len(line) > 5:
		return line, nil
	default:
		return "", errors.New("empty or short passwords are not allowed")
	}
}

func confirmPassword(in *internal.Input, password string) error {
	in.WriteLine("Confirm password:")
	switch read, line, err := in.ReadSensitiveLine(); {
	case err != nil:
		return err
	case !read || line != password:
		return errors.New("passwords mismatch")
	default:
		return nil
	}
}
