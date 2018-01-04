package betareduce

import (
	"encoding/json"
	"errors"
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

const (
	_ = iota
	MSG_CONNECT
	MSG_PUT
	MSG_GET
	MSG_DELETE
)

func REQ_PORT(x int) int { return x }
func SUB_PORT(x int) int { return x + 1 }
func PUB_PORT(x int) int { return x + 2 }

var Msgtypes = map[int]string{
	MSG_CONNECT: "MSG_CONNECT",
	MSG_PUT:     "MSG_PUT",
	MSG_GET:     "MSG_GET",
	MSG_DELETE:  "MSG_DELETE",
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

// Wait for pubsub data (from other betareduce servers )
func recvLoop() {

	// TODO: connect to other replicas here
	p_out("In recvLoop")

	for {
		msg := recv(subSock)
		// TODO: lock here
		p_out("Recv msg %q\n", msg.S)
		// TODO: unlock here
	}
}

// Wait for requests from clients
func repLoop() {

	for {
		p_out("In repLoop")

		msg := recv(repSock)
		// TODO: lock here (on a different lock than recvLoop)
		p_out("Recv msg %q\n", msg.S)
		// TODO: unlock here
	}
}

// Greg TODO
func send(sock *zmq.Socket, m *Msg) {
	fmt.Println("TODO: send msg: " + m.S)
	m.From = me.pubAddr
	s, _ := json.Marshal(m)
	p_out("SEND %q (%d - %d): seq %d, len %d\n", msgtypes[m.Mtype], m.From, m.To, m.Seqnum, len(s))
	bytes, err := sock.SendBytes(s, 0)
	if (err != nil) || (bytes != len(s)) {
		p_err("SEND error, %d bytes, err: %v\n", bytes, err)
		return errors.New("SEND error")
	}
	sends++
	msgTypeSends[m.Mtype]++
	sendBytes += len(s)
	return nil

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
