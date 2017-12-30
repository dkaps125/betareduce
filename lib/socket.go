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
)

func REQ_PORT(x int) int { return x }
func SUB_PORT(x int) int { return x + 1 }

var msgtypes = map[int]string{
	MSG_CONNECT: "MSG_CONNECT",
	MSG_PUT:     "MSG_PUT",
	MSG_GET:     "MSG_GET",
}

type Msg struct {
	S string
	// TODO: put other message info here
}

func InitServer() {
	subsock, err := zmq.NewSocket(zmq.SUB)
	p_dieif(err != nil, "Bad SUB sock, %v\n", err)

	subsock.SetSubscribe("")
}

func ConnectToReplicaReqsock(r Replica) {

	s := fmt.Sprintf("tcp://%s:%d", r.address, REQ_PORT(r.port))
	p_out("connect to " + s)
	r.reqSock, _ = zmq.NewSocket(zmq.REQ)
	if err := r.reqSock.Connect(s); err != nil {
		p_err("Error: cannot connect to server %q, %v\n", s, err)
	}
}

func connectToReplicaSubsock(r Replica) {
	s := fmt.Sprintf("tcp://%s:%d", r.address, SUB_PORT(r.port))
	p_out("connect to " + s)
	if err := subSock.Connect(s); err != nil {
		p_err("Error: cannot connect to server %q, %v\n", s, err)
	}
}

// Wait for pubsub data (from other betareduce servers )
func recvLoop() {

	// TODO: connect to other replicas here
	p_out("In recvLoop")

	for {
		msg := recv(pubSock)
		// TODO: lock here
		p_out("Recv msg %q\n", msg.S)
		// TODO: unlock here
	}
}

// Wait for requests from clients
func repLoop() {

	for {
		p_out("In repLoop")

		msg := recv(pubSock)
		// TODO: lock here (on a different lock than recvLoop)
		p_out("Recv msg %q\n", msg.S)
		// TODO: unlock here
	}
}

func send() {

}

func recv(sock *zmq.Socket) *Msg {
	return nil
}
