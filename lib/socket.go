package betareduce

import (
	"fmt"

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
)

func REQ_PORT(x int) int { return x }
func SUB_PORT(x int) int { return x + 1 }

var Msgtypes = map[int]string{
	MSG_CONNECT:         "MSG_CONNECT",
	MSG_PUT:             "MSG_PUT",
	MSG_PUT_RESPONSE:    "MSG_PUT_RESPONSE",
	MSG_GET:             "MSG_GET",
	MSG_GET_RESPONSE:    "MSG_GET_RESPONSE",
	MSG_DELETE:          "MSG_DELETE",
	MSG_DELETE_RESPONSE: "MSG_DELETE_RESPONSE",
}

type Msg struct {
	S       string
	MsgType int
	// TODO: put other message info here
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

// Greg TODO
func send(sock *zmq.Socket, m *Msg) {
	fmt.Println("TODO: send msg: " + m.S)

}

// Greg TODO
func recv(sock *zmq.Socket) *Msg {
	return nil
}

// Client code blocking send/recv, perhaps move later
func (r *Replica) SendRecv(m *Msg) *Msg {
	send(r.reqSock, m)
	return recv(r.reqSock)
}
