package encryption

import (
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
)

type Compressor interface {
	Compress(s settings.Settings, password string) error
}

type Extractor interface {
	Extract(s settings.Settings, password string) error
}

func Encrypt(in internal.Input) error {
	return encrypt(in, newZipper(&in))
}

func Decrypt(in internal.Input) error {
	return decrypt(in, newZipper(&in))
}
