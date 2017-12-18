package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "github.com/mattn/go-getopt"
)

// add command interpreter

func main() {
	var c, port int
	var address string

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

	for {
		fmt.Print("β> ")
		input, err := reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if err != nil {
			break
		}
	}
}
