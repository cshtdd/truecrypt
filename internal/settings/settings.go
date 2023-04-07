// Settings management package

package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tddapps.com/truecrypt/internal/paths"
)

// Reads the default settings path.
// Supports environment variable overrides.
func DefaultSettingsPath() string {
	if override, found := os.LookupEnv(envSettings); found {
		return override
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic("Cannot read PWD")
	}

	return filepath.Join(pwd, ".settings.json")
}

// Reads the default decrypted folder.
// Supports environment variable overrides.
func DefaultDecryptedFolder() string {
	if override, found := os.LookupEnv(envDecrypted); found {
		return override
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Cannot read $HOME")
	}

	return filepath.Join(homeDir, "decrypted_folder/t")
}

type Settings struct {
	DecryptedFolder paths.Path
	EncryptedFile   paths.Path
}

func LoadFrom(path paths.Path) (Settings, error) {
	result := Settings{}
	bytes, err := path.Read()
	if err != nil {
		return result, err
	}

	json.Unmarshal(bytes, &result)
	return result, nil
}

func (s Settings) Save(path paths.Path) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := path.Write(bytes); err != nil {
		return err
	}

	return nil
}
