package main

import (
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
	flag.StringVar(&flagSettingsPath, "settings", settings.DefaultSettingsPath(), "[Optional] Settings file path")
	flag.BoolVar(&flagPause, "pause", false, "[Optional] Pause at the end")

	flag.Parse()

	// Create Program Input
	input = &internal.Input{
		IO:           internal.IO{Reader: os.Stdin, Writer: os.Stdout},
		SettingsPath: paths.Path(flagSettingsPath),
	}
}

type ExitCode int

const Success ExitCode = 0

type Program struct {
	Input     internal.Input
	Name      string
	ErrorCode ExitCode
}

func NewProgram(name string, errorCode ExitCode) *Program {
	return &Program{
		Input:     *input,
		Name:      name,
		ErrorCode: errorCode,
	}
}

func (p *Program) Run(program func(internal.Input) error) ExitCode {
	p.Input.WriteLine(p.Name)

	if err := program(p.Input); err != nil {
		p.Input.WriteLine(fmt.Sprintf("%s error %s", p.Name, err))
		return p.ErrorCode
	}

	return Success
}

func main() {
	exitCode := Success

	switch {
	case flagSetup:
		exitCode = NewProgram("Setup", 2).Run(setup.Run)
	case flagDecrypt:
		exitCode = NewProgram("Decrypt", 3).Run(encryption.Decrypt)
	case flagEncrypt:
		exitCode = NewProgram("Encrypt", 4).Run(encryption.Encrypt)
	case flagClean:
		exitCode = NewProgram("Clean", 5).Run(clean.Run)
	case flagCleanSettings:
		exitCode = NewProgram("Clean Settings", 6).Run(clean.CleanSettings)
	default:
		input.WriteLine("Action missing")
		flag.Usage()
		exitCode = 1
	}

	if flagPause {
		input.Pause()
	}

	if exitCode != Success {
		os.Exit(int(exitCode))
	}
}
