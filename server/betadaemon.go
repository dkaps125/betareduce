package main

import (
	"betareduce/lib"
	"os"
	"strconv"

	. "github.com/mattn/go-getopt"
)

func main() {
	var c, port int
	debug := false

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
			debug = true
		default:
			os.Exit(1)
		}
	}

	betareduce.Init(port, debug)
}
