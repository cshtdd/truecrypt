package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal/settings"
)

type fakeSettings struct {
	err error
	s   settings.Settings
}

func newFakeSettings(s settings.Settings) *fakeSettings {
	return &fakeSettings{s: s}
}

func newFakeSettingsWithError() *fakeSettings {
	return &fakeSettings{err: errors.New("fake settings loading error")}
}

func (f *fakeSettings) load() (settings.Settings, error) {
	return f.s, f.err
}

func (f *fakeSettings) Err() error {
	return f.err
}
