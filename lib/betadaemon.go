package betareduce

import (
	"betareduce/lib/store"
	"os"
	"strconv"

	. "github.com/mattn/go-getopt"
)

func main() {
	var c, port int

	for {
		if c = Getopt("p:"); c == EOF {
			break
		}

		switch c {
		case 'p':
			port, _ = strconv.Atoi(OptArg)
		default:
			os.Exit(1)
		}
	}

	store.Init(port)
}
