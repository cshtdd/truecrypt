// Settings management package

package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	DecryptedFolder string
	EncryptedFile   string
}

func LoadFrom(path string) (Settings, error) {
	result := Settings{}
	bytes, e := os.ReadFile(path)
	if e != nil {
		return result, e
	}

	json.Unmarshal(bytes, &result)
	return result, nil
}

func (s Settings) Save(path string) error {
	bytes, e := json.Marshal(s)
	if e != nil {
		return e
	}

	const userRwOthersR = 0644
	if e := os.WriteFile(path, bytes, userRwOthersR); e != nil {
		return e
	}

	return nil
}

func (s Settings) IsValidEncryptedFile() (bool, error) {
	if _, err := os.Stat(s.EncryptedFile); err != nil {
		return false, err
	}
	return true, nil
}
