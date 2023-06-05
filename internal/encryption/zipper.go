package encryption

import (
	"os/exec"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
)

// zipper Default Compressor & Extractor
type zipper internal.IO

func newZipper(in *internal.Input) zipper {
	return zipper(in.IO)
}

func (z zipper) Compress(s settings.Settings, password string) error {
	if _, err := exec.Command("zip", "-v").Output(); err != nil { // TODO display this output
		return err
	}

	tmp, err := paths.CreateTempZipFile()
	if err != nil {
		return err
	}

	if err := tmp.Delete(); err != nil {
		return err
	}

	cmd := exec.Command("zip", "-er", "-P", password, tmp.FullPath(), ".")
	cmd.Dir = s.DecryptedFolder.FullPath()
	if _, err := cmd.Output(); err != nil { // TODO: display this output
		return err
	}

	return tmp.Move(s.EncryptedFile)
}

func (z zipper) Extract(s settings.Settings, password string) error {
	if _, err := exec.Command("unzip", "-v").Output(); err != nil { // TODO display this output
		return err
	}

	if err := s.DecryptedFolder.Create(); err != nil {
		return err
	}

	cmd := exec.Command("unzip", "-P", password, "-n", s.EncryptedFile.FullPath(), "-d", s.DecryptedFolder.FullPath())
	_, err := cmd.Output() // TODO: display this output
	return err
}
