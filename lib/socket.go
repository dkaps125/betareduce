package betareduce

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall"

	zmq "github.com/pebbe/zmq4"
)

// ========================================================================== //

const (
	_ = iota
	MSG_CONNECT
	MSG_PUT
	MSG_GET
	MSG_DELETE
	MSG_PUT_RESPONSE
	MSG_GET_RESPONSE
	MSG_DELETE_RESPONSE
	MSG_STATUS
	MSG_SUBSCRIBE
	MSG_SUBSCRIBE_RESPONSE
	MSG_SUBSCRIBE_FLOOD
)

var Msgtypes = map[int]string{
	MSG_CONNECT:            "MSG_CONNECT",
	MSG_PUT:                "MSG_PUT",
	MSG_PUT_RESPONSE:       "MSG_PUT_RESPONSE",
	MSG_GET:                "MSG_GET",
	MSG_GET_RESPONSE:       "MSG_GET_RESPONSE",
	MSG_DELETE:             "MSG_DELETE",
	MSG_DELETE_RESPONSE:    "MSG_DELETE_RESPONSE",
	MSG_STATUS:             "MSG_STATUS",
	MSG_SUBSCRIBE:          "MSG_SUBSCRIBE",
	MSG_SUBSCRIBE_RESPONSE: "MSG_SUBSCRIBE_RESPONSE",
	MSG_SUBSCRIBE_FLOOD:    "MSG_SUBSCRIBE_FLOOD",
}

func CLIENT_PORT(x int) int { return x }
func REP_PORT(x int) int    { return x + 1 }
func BOOT_PORT(x int) int   { return x + 2 }

type Msg struct {
	S       string
	Key     string
	Value   []byte
	Type    string
	MsgType int
	// TODO: put other message info here
	Status int
	From   string
	To     string
	Subs   []Replica
}

// ========================================================================== //

func send(sock *zmq.Socket, m *Msg) error {
	//m.From = me.pubAddr
	s, _ := json.Marshal(m)
	P_out("SEND %q (%q - %q): seq %d, len %d\n", Msgtypes[m.MsgType], m.From, m.To, len(s))

	bytes, err := sock.SendBytes(s, 0)
	if (err != nil) || (bytes != len(s)) {
		P_err("SEND error, %d bytes, err: %v\n", bytes, err)
		return errors.New("SEND error")
	}
	return nil

}

func recv(sock *zmq.Socket) *Msg {
	var str string
	var err error
	var flags zmq.Flag

	// do we need to do something with this?
	flags = 0

	m := new(Msg)
	for {
		str, err = sock.Recv(flags)
		if err == nil {
			break
		}

		if err.Error() == "interrupted system call" {
			P_out("--System call interrupted--")
		} else if err.Error() == "resource temporarily unavailable" {
			return nil
		} else {
			P_die("recv err: %q (%v), \n", err.Error(), syscall.EINTR)
		}
	}
	if err := json.Unmarshal([]byte(str), m); err != nil {
		P_err("ERROR unmarshaling message %q\n", string(str))
	}
	P_out("\n\tRECV %q (%q - %q): len %d\n\n", Msgtypes[m.MsgType], m.From, m.To, len(str))
	return m
}

// ConnectToReplicaReqsock is a function called by the client to
func ConnectToReplicaReqsock(address string, port int) Replica {
	r := Replica{
		address: address,
		port:    port,
	}

	s := fmt.Sprintf("tcp://%s:%d", r.address, CLIENT_PORT(r.port))
	P_out("connect to " + s)
	r.reqSock, _ = zmq.NewSocket(zmq.REQ)
	if err := r.reqSock.Connect(s); err != nil {
		P_err("Error: cannot connect to server %q, %v\n", s, err)
		panic(err)
	}

	return r
}
