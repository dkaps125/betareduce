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
		if c = Getopt("da:p:"); c == EOF {
			break
		}
		switch c {
		case 'a':
			address = OptArg
		case 'p':
			port, _ = strconv.Atoi(OptArg)
		case 'd':
			Debug = true
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

		if len(op) > 1 {
			switch op[0] {
			case "put":

				if len(op) < 3 {
					break
				}
				outboundMsg := &Msg{
					MsgType: MSG_PUT,
					Key:     op[2],
					Type:    op[1],
					Value:   String{Value: strings.Join(op[3:], " ")}.Serialize(),
				}

				//fmt.Printf("Sending %s, %s\n", outboundMsg.Key, outboundMsg.Value)

				replyMsg := replica.SendRecv(outboundMsg)
				fmt.Printf("PUT %s, %v\n", replyMsg.Key, GetValue(replyMsg.Value, op[1]))
				break
			case "get":
				outboundMsg := &Msg{
					MsgType: MSG_GET,
					Key:     op[1],
				}
				replyMsg := replica.SendRecv(outboundMsg)
				//TODO: change types so that they are specified in serialization of value
				fmt.Printf("GET %s, %v\n", replyMsg.Key, GetValue(replyMsg.Value, "String"))
				break
			default:
				fmt.Println("Command not recognized")
				break
			}
		}

	}
}

// add functions to send commands to replicas
