package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cshtdd/truecrypt/internal"
	"github.com/cshtdd/truecrypt/internal/clean"
	"github.com/cshtdd/truecrypt/internal/encryption"
	"github.com/cshtdd/truecrypt/internal/paths"
	"github.com/cshtdd/truecrypt/internal/settings"
	"github.com/cshtdd/truecrypt/internal/setup"
	"os"
)

var flagSetup bool
var flagEncrypt bool
var flagDecrypt bool
var flagClean bool
var flagCleanSettings bool
var flagSettingsPath string
var flagPause bool

var input *internal.Input

func init() {
	// Commands
	flag.BoolVar(&flagSetup, "setup", false, "Setup settings")
	flag.BoolVar(&flagEncrypt, "encrypt", false, "Encrypts decrypted folder")
	flag.BoolVar(&flagDecrypt, "decrypt", false, "Decrypts encrypted folder")
	flag.BoolVar(&flagClean, "clean", false, "Deletes the decrypted folder")
	flag.BoolVar(&flagCleanSettings, "cleanSettings", false, "Deletes the settings")

	// Other config
	flag.StringVar(&flagSettingsPath, "settings", settings.DefaultSettingsPath().String(), "[Optional] Settings file path")
	flag.BoolVar(&flagPause, "pause", false, "[Optional] Pause at the end")

	flag.Parse()

	// Create Program Input
	input = &internal.Input{
		IO:           internal.IO{Reader: os.Stdin, Writer: os.Stdout},
		SettingsPath: paths.FilePath(flagSettingsPath),
	}
}

type ExitCode int

type Program struct {
	Input     *internal.Input
	Name      string
	ErrorCode ExitCode
}

func NewProgram(name string, errorCode ExitCode) *Program {
	return &Program{
		Input:     input,
		Name:      name,
		ErrorCode: errorCode,
	}
}

func (p *Program) Run(program func(*internal.Input) error) {
	p.Input.WriteLine(p.Name)

	err := program(p.Input)
	if err != nil {
		p.Input.WriteLine(fmt.Sprintf("%s error %s", p.Name, err))
	}

	if flagPause {
		input.Pause()
	}

	if err != nil {
		os.Exit(int(p.ErrorCode))
	}
}

func main() {
	switch {
	case flagSetup:
		NewProgram("Setup", 2).Run(setup.Run)
	case flagDecrypt:
		NewProgram("Decrypt", 3).Run(encryption.Decrypt)
	case flagEncrypt:
		NewProgram("Encrypt", 4).Run(encryption.Encrypt)
	case flagClean:
		NewProgram("Clean", 5).Run(clean.Run)
	case flagCleanSettings:
		NewProgram("Clean Settings", 6).Run(clean.Settings)
	default:
		NewProgram("Action Missing", 1).Run(func(i *internal.Input) error {
			flag.Usage()
			return errors.New("empty arguments")
		})
	}
}
