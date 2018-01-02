// leave this as main
package main

import (
	betareduce "betareduce/lib"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "github.com/mattn/go-getopt"
)

// add command interpreter

func main() {
	var c int
	var address string
	var port int

	address = "127.0.0.1"
	port = 8300

	for {
		if c = Getopt("a:"); c == EOF {
			break
		}
		switch c {
		case 'a':
			address = OptArg
		case 'p':
			port, _ = strconv.Atoi(OptArg)
		default:
			println("usage: betareduce.go [-a address]", c)
			os.Exit(1)
		}
	}

	// connect to server here

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Connecting to tcp://" + address + ":" + string(port))
	replica := betareduce.ConnectToReplicaReqsock(address, port)

	for {
		fmt.Print("β> ")
		input, err := reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if err != nil {
			break
		}
	}
}

// add functions to send commands to replicas
