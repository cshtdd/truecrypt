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
}

func main() {
	const Success = 0
	exitCode := Success

	i := internal.Input{
		IO:           internal.IO{Reader: os.Stdin, Writer: os.Stdout},
		SettingsPath: paths.Path(flagSettingsPath),
	}

	switch {
	case flagSetup:
		program := "Setup"
		i.WriteLine(program)
		if err := setup.Run(i); err != nil {
			i.WriteLine(fmt.Sprintf("%s error %s", program, err))
			exitCode = 2
		}
	case flagDecrypt:
		program := "Decrypt"
		i.WriteLine(program)
		if err := encryption.Decrypt(i); err != nil {
			i.WriteLine(fmt.Sprintf("%s error %s", program, err))
			exitCode = 3
		}
	case flagEncrypt:
		program := "Encrypt"
		i.WriteLine(program)
		if err := encryption.Encrypt(i); err != nil {
			i.WriteLine(fmt.Sprintf("%s error %s", program, err))
			exitCode = 4
		}
	case flagClean:
		program := "Clean"
		i.WriteLine(program)
		if err := clean.Run(i); err != nil {
			i.WriteLine(fmt.Sprintf("%s error %s", program, err))
			exitCode = 5
		}
	case flagCleanSettings:
		program := "Clean Settings"
		i.WriteLine(program)
		if err := clean.CleanSettings(i); err != nil {
			i.WriteLine(fmt.Sprintf("%s error %s", program, err))
			exitCode = 6
		}
	default:
		i.WriteLine("Action missing")
		flag.Usage()
		exitCode = 1
	}

	if flagPause {
		i.Pause()
	}

	if exitCode != Success {
		os.Exit(exitCode)
	}
}
