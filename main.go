package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/clean"
	"tddapps.com/truecrypt/internal/encryption"
	"tddapps.com/truecrypt/internal/paths"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
)

var flagSetup bool
var flagEncrypt bool
var flagDecrypt bool
var flagClean bool
var flagCleanSettings bool
var flagTest bool
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
	flag.BoolVar(&flagTest, "test", false, "Runs test program")

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
	case flagTest:
		NewProgram("Test", 7).Run(func(i *internal.Input) error {
			i.WriteLine("Enter secret:")
			switch read, password, err := i.ReadSensitiveLine(); {
			case err != nil:
				return err
			case !read:
				return errors.New("empty input")
			default:
				i.WriteLine(fmt.Sprintf("Thanks for your input p: %s", password))
				return nil
			}
		})
	default:
		NewProgram("Action Missing", 1).Run(func(i *internal.Input) error {
			flag.Usage()
			return errors.New("empty arguments")
		})
	}
}
