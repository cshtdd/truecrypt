package settings_test

import (
	"fmt"
	"os"
	"testing"

	"tddapps.com/truecrypt/internal/settings"
)

func TestDefaultSettingsPath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Cannot read pwd %s", err)
	}

	expected := fmt.Sprintf("%s/.settings.json", pwd)
	path := settings.DefaultSettingsPath()
	if path != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", path, expected)
	}
}

func TestOverridesPathFromEnv(t *testing.T) {
	expected := "aaa/bbb.json"
	os.Setenv("TC_SETTINGS", expected)
	defer os.Unsetenv("TC_SETTINGS")

	path := settings.DefaultSettingsPath()
	if path != expected {
		t.Errorf("DefaultSettingsPath() = %s, want %s", path, expected)
	}
}
