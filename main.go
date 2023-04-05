package main

import (
	"flag"
	"fmt"
	"os"

	"tddapps.com/truecrypt/internal/settings"
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
	switch {
	case flagSetup:
		fmt.Println("Setup")
		// TODO: implement this
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
		os.Exit(1)
	}
}
