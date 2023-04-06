package main

import (
	"flag"
	"fmt"
	"os"

	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/settings"
	"tddapps.com/truecrypt/internal/setup"
)

var flagSetup bool
var flagEncrypt bool
var flagDecrypt bool
var flagClean bool
var flagSettingsPath string

func init() {
	// Commands
	flag.BoolVar(&flagSetup, "setup", false, "Setup settings")
	flag.BoolVar(&flagEncrypt, "encrypt", false, "Encrypts decrypted folder")
	flag.BoolVar(&flagDecrypt, "decrypt", false, "Decrypts encrypted folder")
	flag.BoolVar(&flagClean, "clean", false, "Deletes the decrypted folder")

	// Other config
	flag.StringVar(&flagSettingsPath, "settings", settings.DefaultSettingsPath(), "[Optional] Settings file path")

	flag.Parse()
}

func main() {
	const Success = 0
	exitCode := Success

	io := internal.IO{Reader: os.Stdin, Writer: os.Stdout}

	switch {
	case flagSetup:
		fmt.Println("Setup")
		i := setup.Input{IO: io, SettingsPath: flagSettingsPath}
		if e := setup.Run(i); e != nil {
			fmt.Println("Setup error", e)
			exitCode = 1
		}
	case flagDecrypt:
		fmt.Println("Decrypt")
		// TODO: implement this
	case flagEncrypt:
		fmt.Println("Encrypt")
		// TODO: implement this
	case flagClean:
		fmt.Println("Clean")
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
