package settings

import (
	"os"
	"tddapps.com/truecrypt/internal/paths"
)

// DefaultSettingsPath Reads the default settings path.
// Supports environment variable overrides.
func DefaultSettingsPath() paths.FilePath {
	if override, found := os.LookupEnv(envSettings); found {
		return paths.FilePath(override)
	}

	return "~/.config/truecrypt/settings.json"
}

// DefaultDecryptedFolder Reads the default decrypted folder.
// Supports environment variable overrides.
func DefaultDecryptedFolder() paths.DirPath {
	if override, found := os.LookupEnv(envDecrypted); found {
		return paths.DirPath(override)
	}

	return "~/decrypted_folder/t"
}
