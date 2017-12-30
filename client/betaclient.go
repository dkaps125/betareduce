// leave this as main
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

	address = "127.0.0.1:8300"

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
	}

	// connect to server here

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Connecting to tcp://" + address)
	betareduce.ConnectToReplica(address)

	for {
		fmt.Print("β> ")
		input, err := reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if err != nil {
			break
		}
	}
}
