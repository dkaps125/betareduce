package main

import (
	"os"
	"strconv"

	. "github.com/mattn/go-getopt"

	lib "betareduce/lib"
)

func main() {
	var c, port int

	// uh here's a default port
	port = 8300

	for {
		if c = Getopt("dp:"); c == EOF {
			break
		}

		switch c {
		case 'p':
			port, _ = strconv.Atoi(OptArg)
		case 'd':
			lib.Debug = true
		default:
			os.Exit(1)
		}
	}

	Run(port)
}
