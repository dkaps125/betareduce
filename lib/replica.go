package betareduce

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

var err error

// Replica is a representation of a single server in a replica set
// for a distributed key-value store
type Replica struct {
	address  string
	port     int
	reqSock  *zmq.Socket
	repSock  *zmq.Socket
	subSock  *zmq.Socket
	pubSock  *zmq.Socket
	bootSock *zmq.Socket
	isLive   bool
	pubAddr  string

	repList []Replica
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

	s := fmt.Sprintf("tcp://*:%d", REP_PORT(r.port))
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

	s = fmt.Sprintf("tcp://*:%d", CLIENT_PORT(r.port))
	if err = r.repSock.Bind(s); err != nil {
		P_die("Error binding req/rep socket %q (%v)\n", s, err)
	}

	P_out("request bound to %q\n", s)
}

// Client code blocking send/recv, perhaps move later
func (r *Replica) SendRecvToReplica(m *Msg) *Msg {
	inClientSend()
	P_out("Sending")
	send(r.reqSock, m)
	outClientSend()
	inClientRecv()
	defer outClientRecv()
	P_out("Exit send")
	return recv(r.reqSock)
}

func (r *Replica) RecvClient() *Msg {
	inClientRecv()
	defer outClientRecv()
	return recv(r.repSock)
}

func (r *Replica) SendClient(m *Msg) {
	inClientSend()
	defer outClientSend()
	send(r.repSock, m)
}

func (r *Replica) RecvRep() *Msg {
	inRepRecv()
	defer outRepRecv()
	return recv(r.subSock)
}

func (r *Replica) SendRep(m *Msg) {
	inRepSend()
	defer outRepSend()
	m.From = fmt.Sprintf("%s:%d", r.address, r.port)
	send(r.pubSock, m)
}

func (r *Replica) SendRecvBoot(m *Msg) *Msg {
	P_out("Sending")
	r.SendBoot(m)
	P_out("Exit send")
	return r.RecvBoot()
}

func (r *Replica) RecvBoot() *Msg {
	inBootRecv()
	defer outBootRecv()
	return recv(r.bootSock)
}

func (r *Replica) SendBoot(m *Msg) {
	inBootSend()
	defer outBootSend()
	m.From = fmt.Sprintf("%s:%d", r.address, r.port)
	send(r.bootSock, m)
}

func (r *Replica) Subscribe(address string, port int) {
	s := fmt.Sprintf("tcp://%s:%d", address, REP_PORT(port))
	P_out("connect to " + s)
	if err = r.subSock.Connect(s); err != nil {
		P_err("Error: cannot connect to server %q, %v\n", s, err)
	}
}

func (r *Replica) BootstrapFromReplica(bootstrapAddress string) {
	if bootstrapAddress != "" {
		address, port := GetAddrComponents(bootstrapAddress)

		s := fmt.Sprintf("tcp://%s:%d", address, BOOT_PORT(port))
		P_out("connect to " + s)
		r.bootSock, err = zmq.NewSocket(zmq.REQ)
		if err = r.bootSock.Connect(s); err != nil {
			P_err("Error: cannot connect to server %q, %v\n", s, err)
			panic(err)
		}

		m := &Msg{
			MsgType: MSG_SUBSCRIBE,
		}

		m = r.SendRecvBoot(m)

		if m.Status == 0 {
			for _, v := range m.Subs {
				r.Subscribe(v.address, v.port)
			}

			r.Subscribe(address, port)
		} else {
			panic(err)
		}

		// TODO: populate KVS with everything from boostrap node
	}

	r.bootSock, err = zmq.NewSocket(zmq.REP)

	s := fmt.Sprintf("tcp://*:%d", BOOT_PORT(r.port))
	if err = r.bootSock.Bind(s); err != nil {
		P_die("Error binding req/rep socket %q (%v)\n", s, err)
	}

	P_out("boot bound to %q\n", s)
}

func (r *Replica) HandleBootstrap(m *Msg) {
	flood := &Msg{
		MsgType: MSG_SUBSCRIBE_FLOOD,
		Status:  0,
		S:       m.From,
	}

	r.SendRep(flood)

	// TODO: add sender to repList, and subscribe
	address, port := GetAddrComponents(m.From)
	r.repList = append(r.repList, Replica{address: m.From})

	r.Subscribe(address, port)

	resp := &Msg{
		MsgType: MSG_SUBSCRIBE_RESPONSE,
		Status:  0,
		Subs:    r.repList,
	}

	r.SendBoot(resp)

}
