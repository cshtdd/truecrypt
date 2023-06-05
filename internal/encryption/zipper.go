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
	if err := z.Run(exec.Command("zip", "-v")); err != nil {
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
	if err := z.Run(cmd); err != nil {
		return err
	}

	return tmp.Move(s.EncryptedFile)
}

func (z zipper) Extract(s settings.Settings, password string) error {
	if err := z.Run(exec.Command("unzip", "-v")); err != nil {
		return err
	}

	if err := s.DecryptedFolder.Create(); err != nil {
		return err
	}

	cmd := exec.Command("unzip", "-P", password, "-n", s.EncryptedFile.FullPath(), "-d", s.DecryptedFolder.FullPath())
	return z.Run(cmd)
}

func (z zipper) Run(cmd *exec.Cmd) error {
	o, err := cmd.Output()
	z.Writer.Write(o)
	return err
}
