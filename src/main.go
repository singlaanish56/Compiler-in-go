package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/singlaanish56/Compiler-in-go/repl"
)


func main(){

	user, err := user.Current()
	if err != nil{
		panic(err)
	}

	fmt.Printf("Hello %s, welcome to the REPL!\n", user.Username)
	fmt.Println("Go ahead type something")
	repl.Start(os.Stdin, os.Stdout)
}