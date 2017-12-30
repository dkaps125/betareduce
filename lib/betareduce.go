package betareduce

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

var (
	store     KVS
	storeLock chan (int)

	pubSock, subSock, repSock *zmq.Socket
	me                        *Replica

	debug = false
)

type Replica struct {
	address string
	port    int
	reqSock *zmq.Socket
	isLive  bool
}

// ========================================================================== //

func in() {
	storeLock <- 1
}

func out() {
	<-storeLock
}

// ========================================================================== //

// Init initializes a key-value store and binds to a port
func Init(port int, _debug bool) {
	store = NewKVS()
	storeLock = make(chan (int), 1)

	debug = _debug

	me = &Replica{
		address: "127.0.0.1",
		port:    port,
		isLive:  true,
	}

	// bind pub/sub
	ps, err := zmq.NewSocket(zmq.PUB)

	pubSock = ps
	if err != nil {
		p_out("skt create err %v\n", err)
	}
	s := fmt.Sprintf("tcp://*:%d", SUB_PORT(port))
	err = pubSock.Bind(s)
	if err != nil {
		p_out("Error binding pub/sub socket %q (%v)\n", s, err)
	}
	p_out("pubsub bound to %q\n", s)

	// bind req/rep
	rs, _ := zmq.NewSocket(zmq.REP)

	repSock = rs

	s = fmt.Sprintf("tcp://*:%d", REQ_PORT(port))
	if err := repSock.Bind(s); err != nil {
		p_die("Error binding req/rep socket %q (%v)\n", s, err)
	}

	p_out("request bound to %q\n", s)

	go recvLoop()
	go repLoop()
}

// ========================================================================== //
// Key value function wrappers

func put(key string, value Value) {
	in()
	store.Put(key, value)
	out()
}

func get(key string) (Value, error) {
	in()
	s, e := store.Get(key)
	out()

	if e == nil {
		return nil, EKEYNF
	}

	return s, nil
}
