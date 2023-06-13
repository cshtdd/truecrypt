package settings_test

import (
	"testing"

	"github.com/cshtdd/truecrypt/internal/settings"
	"github.com/cshtdd/truecrypt/internal/test/helpers"
)

func TestSettings_Serialization(t *testing.T) {
	tmp := helpers.CreateTemp(t)

	s := settings.Settings{DecryptedFolder: "aaa/bbb", EncryptedFile: "ccc.zip"}
	if err := s.Save(tmp); err != nil {
		t.Fatalf("Unexpected error err: %s", err)
	}

	x, err := settings.LoadFrom(tmp)
	if err != nil {
		t.Fatal("Got unexpected error", err)
	}

	if x != s {
		t.Errorf("got %v+, want %v+", x, s)
	}
}
