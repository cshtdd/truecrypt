package main

import (
	"fmt"
	"os"
)

func main() {
	pwd, _ := os.Getwd()
	fmt.Println("Truecrypt ", "pwd:", pwd, "args:", os.Args)
}
