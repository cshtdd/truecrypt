package main

import (
	"flag"
	"fmt"
	"os"
)

var flagSetup bool
var flagEncrypt bool
var flagDecrypt bool
var flagClean bool

func init() {
	flag.BoolVar(&flagSetup, "setup", false, "Setup settings")
	flag.BoolVar(&flagEncrypt, "encrypt", false, "Encrypts decrypted folder")
	flag.BoolVar(&flagDecrypt, "decrypt", false, "Decrypts encrypted folder")
	flag.BoolVar(&flagClean, "clean", false, "Deletes the decrypted folder")
	flag.Parse()
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot read PWD")
		os.Exit(1)
	}
	fmt.Println("Truecrypt", "pwd:", pwd, "flags:", flag.Args())

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
