package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/mannyjimen/Monkey-Compiler/repl"
)

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %+v, welcome to Monkey!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
