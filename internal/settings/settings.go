// Settings management package

package settings

import (
	"encoding/json"
	"os"

	"tddapps.com/truecrypt/internal/paths"
)

// Reads the default settings path.
// Supports environment variable overrides.
func DefaultSettingsPath() paths.FilePath {
	if override, found := os.LookupEnv(envSettings); found {
		return paths.FilePath(override)
	}

	return "~/.config/truecrypt/settings.json"
}

// Reads the default decrypted folder.
// Supports environment variable overrides.
func DefaultDecryptedFolder() paths.DirPath {
	if override, found := os.LookupEnv(envDecrypted); found {
		return paths.DirPath(override)
	}

	return "~/decrypted_folder/t"
}

type Settings struct {
	DecryptedFolder paths.DirPath
	EncryptedFile   paths.FilePath // TODO: change this to ZipPath
}

func LoadFrom(p paths.FilePath) (Settings, error) {
	result := Settings{}
	bytes, err := p.Read()
	if err != nil {
		return result, err
	}

	json.Unmarshal(bytes, &result)
	return result, nil
}

func (s Settings) Save(p paths.FilePath) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := p.Write(bytes); err != nil {
		return err
	}

	return nil
}
