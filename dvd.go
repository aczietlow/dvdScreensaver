package main

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	width, height, err := term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	// poor man's kill switch
	//c := make([]byte, 1)
	//os.Stdin.Read(c)
	//switch c[0] {
	//case 'q':
	//	err = term.Restore(int(os.Stdin.Fd()), oldState)
	//	os.Exit(0)
	//
	//}

	terminal := term.NewTerminal(os.Stdin, "")
	fmt.Println("user name?")
	line, err := terminal.ReadLine()
	if err != nil {
		panic(err)
	}

	fmt.Printf("line:%v\n", line)

	fmt.Printf("unix File Descriptor:%v\n", int(os.Stdin.Fd()))
	fmt.Printf("x:%v,y:%v\n", width, height)
	fmt.Println("end of terminal")
}
