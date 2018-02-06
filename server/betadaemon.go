package main

import (
	"os"
	"strconv"

	opt "github.com/mattn/go-getopt"

	lib "betareduce/lib"
)

func main() {
	var c, port int
	bootstrap := ""

	// uh here's a default port
	port = 8300

	for {
		if c = opt.Getopt("dp:b:"); c == opt.EOF {
			break
		}

		switch c {
		case 'p':
			port, _ = strconv.Atoi(opt.OptArg)
		case 'd':
			lib.Debug = true
		case 'b':
			bootstrap = opt.OptArg
		default:
			os.Exit(1)
		}
	}

	Run(port, bootstrap)
}
