// leave this as main
package main

import (
	lib "betareduce/lib"
	"encoding/binary"

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
			lib.Debug = true
		default:
			println("usage: betareduce.go [-a address] [-p port] [-d]", c)
			os.Exit(1)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Connecting to tcp://" + address + ":" + string(port))
	replica := lib.ConnectToReplicaReqsock(address, port)

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

				if len(op) < 4 {
					break
				}
				outboundMsg := &lib.Msg{
					MsgType: lib.MSG_PUT,
					Key:     op[2],
					Value:   getBytes(op[1], strings.Join(op[3:], " ")),
				}

				lib.P_out("Sending %s, %v\n", outboundMsg.Key, outboundMsg.Value)

				replyMsg := replica.SendRecvToReplica(outboundMsg)
				fmt.Printf("PUT %s, %v\n", replyMsg.Key, lib.GetValue(replyMsg.Value))
				break
			case "get":
				outboundMsg := &lib.Msg{
					MsgType: lib.MSG_GET,
					Key:     op[1],
				}
				replyMsg := replica.SendRecvToReplica(outboundMsg)
				//TODO: change types so that they are specified in serialization of value
				if replyMsg.Status == 0 {
					fmt.Printf("GET %s, %v\n", replyMsg.Key, lib.GetValue(replyMsg.Value))
				} else {
					fmt.Print(replyMsg.Key + ": KEY NOT FOUND\n")
				}
				break
			default:
				fmt.Println("Command not recognized")
				break
			}
		}

	}
}

func getBytes(kind string, contents string) []byte {
	var pref byte
	var ret []byte

	switch kind {
	case "String":
		pref = 's'
	case "Int":
		pref = 'i'
	default:
		pref = 'b'
	}

	if pref == 'i' {
		ret = make([]byte, 9)
		ret[0] = 'i'

		tmp, _ := strconv.Atoi(contents)
		binary.BigEndian.PutUint64(ret[1:], uint64(tmp))

		return ret
	}

	return append([]byte{pref}, []byte(contents)...)
}

// add functions to send commands to replicas
