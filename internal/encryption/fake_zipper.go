package encryption

import (
	"errors"
	"tddapps.com/truecrypt/internal/settings"
)

type fakeZipper struct {
	err error
}

func newFakeZipper() *fakeZipper {
	return &fakeZipper{}
}

func newFakeZipperWithError() *fakeZipper {
	return &fakeZipper{
		err: errors.New("fake Zipper error"),
	}
}

func (f *fakeZipper) compress(_ settings.Settings, _ string) error {
	return f.err
}

func (f *fakeZipper) extract(_ settings.Settings, _ string) error {
	return f.err
}

func (f *fakeZipper) Err() error {
	return f.err
}
