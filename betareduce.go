package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/mattn/go-getopt"
)

func main() {
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
