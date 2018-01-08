package betareduce

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

var err error

// Replica is a representation of a single server in a replica set
// for a distributed key-value store
type Replica struct {
	address string
	port    int
	reqSock *zmq.Socket
	repSock *zmq.Socket
	subSock *zmq.Socket
	pubSock *zmq.Socket
	isLive  bool
	pubAddr string
}

// CreateReplica returns a replica with the given address and port
func CreateReplica(address string, port int) *Replica {
	return &Replica{
		address: address,
		port:    port,
	}
}

func (r *Replica) InitReplica() {
	r.pubSock, err = zmq.NewSocket(zmq.PUB)

	if err != nil {
		P_out("skt create err %v\n", err)
	}

	s := fmt.Sprintf("tcp://*:%d", PUB_PORT(r.port))
	err = r.pubSock.Bind(s)

	if err != nil {
		P_out("Error binding pub/sub socket %q (%v)\n", s, err)
	}
	P_out("pub sock bound to %q\n", s)

	r.subSock, err = zmq.NewSocket(zmq.SUB)
	P_dieif(err != nil, "Bad SUB sock, %v\n", err)
	P_out("sub sock initialized to %q\n", r.subSock)

	r.subSock.SetSubscribe("")

	// bind req/rep
	r.repSock, err = zmq.NewSocket(zmq.REP)

	s = fmt.Sprintf("tcp://*:%d", REQ_PORT(r.port))
	if err := r.repSock.Bind(s); err != nil {
		P_die("Error binding req/rep socket %q (%v)\n", s, err)
	}

	P_out("request bound to %q\n", s)
}

// Client code blocking send/recv, perhaps move later
func (r *Replica) SendRecv(m *Msg) *Msg {
	P_out("Sending")
	send(r.reqSock, m)
	P_out("Exit send")
	return recv(r.reqSock)
}

func (r *Replica) RecvClient() *Msg {
	return recv(r.repSock)
}

func (r *Replica) SendClient(m *Msg) {
	send(r.repSock, m)
}
