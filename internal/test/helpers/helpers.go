package helpers

import (
	"os"
	"testing"
)

func CreateTemp(t *testing.T) string {
	tmp, err := os.CreateTemp("", "tc_settings")
	if err != nil {
		t.Fatal("Cannot create tmp file", err)
		return ""
	}
	t.Cleanup(func() {
		os.Remove(tmp.Name())
	})
	return tmp.Name()
}
