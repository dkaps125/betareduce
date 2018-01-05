package betareduce

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall"

	zmq "github.com/pebbe/zmq4"
)

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
)

var Msgtypes = map[int]string{
	MSG_CONNECT:         "MSG_CONNECT",
	MSG_PUT:             "MSG_PUT",
	MSG_PUT_RESPONSE:    "MSG_PUT_RESPONSE",
	MSG_GET:             "MSG_GET",
	MSG_GET_RESPONSE:    "MSG_GET_RESPONSE",
	MSG_DELETE:          "MSG_DELETE",
	MSG_DELETE_RESPONSE: "MSG_DELETE_RESPONSE",
	MSG_STATUS:          "MSG_STATUS",
}

func REQ_PORT(x int) int { return x }
func SUB_PORT(x int) int { return x + 1 }
func PUB_PORT(x int) int { return x + 2 }

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
}

var semSend = make(chan (int), 1)
var semRecv = make(chan (int), 1)

func inSend() {
	semSend <- 1
}

func outSend() {
	<-semSend
}

func inRecv() {
	semRecv <- 1
}

func outRecv() {
	<-semRecv
}

func ConnectToReplicaReqsock(address string, port int) Replica {
	r := Replica{
		address: address,
		port:    port,
	}

	s := fmt.Sprintf("tcp://%s:%d", r.address, REQ_PORT(r.port))
	p_out("connect to " + s)
	r.reqSock, _ = zmq.NewSocket(zmq.REQ)
	if err := r.reqSock.Connect(s); err != nil {
		p_err("Error: cannot connect to server %q, %v\n", s, err)
	}

	return r
}

func connectToReplicaSubsock(r Replica) {
	// add some logic for making sure replica is there, listening

	s := fmt.Sprintf("tcp://%s:%d", r.address, SUB_PORT(r.port))
	p_out("connect to " + s)
	if err := subSock.Connect(s); err != nil {
		p_err("Error: cannot connect to server %q, %v\n", s, err)
	}
}

func send(sock *zmq.Socket, m *Msg) error {
	inSend()
	defer outSend()

	//m.From = me.pubAddr
	s, _ := json.Marshal(m)
	p_out("SEND %q (%q - %q): seq %d, len %d\n", Msgtypes[m.MsgType], m.From, m.To, len(s))

	bytes, err := sock.SendBytes(s, 0)
	if (err != nil) || (bytes != len(s)) {
		p_err("SEND error, %d bytes, err: %v\n", bytes, err)
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
	inRecv()
	defer outRecv()

	m := new(Msg)
	for {
		str, err = sock.Recv(flags)
		if err == nil {
			break
		}

		if err.Error() == "interrupted system call" {
			p_out("--System call interrupted--")
		} else if err.Error() == "resource temporarily unavailable" {
			return nil
		} else {
			p_die("recv err: %q (%v), \n", err.Error(), syscall.EINTR)
		}
	}
	if err := json.Unmarshal([]byte(str), m); err != nil {
		p_err("ERROR unmarshaling message %q\n", string(str))
	}
	p_out("\n\tRECV %q (%q - %q): len %d\n\n", Msgtypes[m.MsgType], m.From, m.To, len(str))
	return m
}

// Client code blocking send/recv, perhaps move later
func (r *Replica) SendRecv(m *Msg) *Msg {
	p_out("Sending")
	send(r.reqSock, m)
	p_out("Exit send")
	return recv(r.reqSock)
}
