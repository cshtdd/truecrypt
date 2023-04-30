package main

import (
	"flag"
	"fmt"
	"os"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/clean"
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

func init() {
	// Commands
	flag.BoolVar(&flagSetup, "setup", false, "Setup settings")
	flag.BoolVar(&flagEncrypt, "encrypt", false, "Encrypts decrypted folder")
	flag.BoolVar(&flagDecrypt, "decrypt", false, "Decrypts encrypted folder")
	flag.BoolVar(&flagClean, "clean", false, "Deletes the decrypted folder")
	flag.BoolVar(&flagCleanSettings, "cleanSettings", false, "Deletes the settings")

	// Other config
	flag.StringVar(&flagSettingsPath, "settings", settings.DefaultSettingsPath(), "[Optional] Settings file path")

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
		fmt.Println("Setup")
		if err := setup.Run(i); err != nil {
			fmt.Println("Setup error", err)
			exitCode = 2
		}
	case flagDecrypt:
		fmt.Println("Decrypt")
		// TODO: implement this
	case flagEncrypt:
		fmt.Println("Encrypt")
		// TODO: implement this
	case flagClean:
		fmt.Println("Clean")
		if err := clean.Run(i); err != nil {
			fmt.Println("Cleanr error", err)
			exitCode = 5
		}
	case flagCleanSettings:
		fmt.Println("Clean Settings")
		// TODO: implement this
	default:
		fmt.Println("Action missing")
		flag.Usage()
		exitCode = 1
	}

	if exitCode != Success {
		os.Exit(exitCode)
	}
}
