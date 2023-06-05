package settings_test

import (
	"testing"

	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/test/helpers"
)

func TestSettings_Serialization(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	s := settings.Settings{DecryptedFolder: "aaa/bbb", EncryptedFile: "ccc.zip"}
	s.Save(tmp)

	x, err := settings.LoadFrom(tmp)
	if err != nil {
		t.Fatal("Got unexpected error", err)
	}

	if x != s {
		t.Errorf("got %v+, want %v+", x, s)
	}
}
