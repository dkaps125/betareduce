package main

import (
	"betareduce/lib"
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/mattn/go-getopt"
)

// add command interpreter

func main() {
	var c int
	var address string

	for {
		if c = Getopt("a:"); c == EOF {
			break
		}
		switch c {
		case 'a':
			address = OptArg
		default:
			println("usage: betareduce.go [-a address]", c)
			os.Exit(1)
		}

		betareduce.Init(address)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("β> ")
		input, err := reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if err != nil {
			break
		}
	}
}
