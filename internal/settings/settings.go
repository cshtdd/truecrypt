// Settings management package

package settings

import (
	"os"
	"path/filepath"
)

// Reads the default settings path.
// Supports environment variable overrides.
func DefaultSettingsPath() string {
	if pathOverride, exists := os.LookupEnv("TC_SETTINGS"); exists {
		return pathOverride
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic("Cannot read PWD")
	}

	return filepath.Join(pwd, ".settings.json")
}
