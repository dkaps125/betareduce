// leave this as main
package main

import (
	. "betareduce/lib"
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

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Connecting to tcp://" + address + ":" + string(port))
	replica := ConnectToReplicaReqsock(address, port)

	for {
		fmt.Print("β> ")
		input, err := reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if err != nil {
			break
		}

		if len(input) == 0 {
			continue
		}

		op := strings.Split(input, " ")

		if len(op) >= 1 {
			switch op[0] {
			case "put":
				outboundMsg := &Msg{
					MsgType: MSG_PUT,
					S:       strings.Join(op[1:], " "),
				}
				replyMsg := replica.SendRecv(outboundMsg)
				fmt.Println(replyMsg)
				break
			case "get":
				outboundMsg := &Msg{
					MsgType: MSG_GET,
					S:       strings.Join(op[1:], " "),
				}
				replyMsg := replica.SendRecv(outboundMsg)
				fmt.Println(replyMsg)
				break
			default:
				fmt.Println("Command not recognized")
			}
		}

	}
}

// add functions to send commands to replicas
