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

	Debug = false
)

type Replica struct {
	address string
	port    int
	reqSock *zmq.Socket
	isLive  bool
	pubAddr string
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
	var err error
	store = NewKVS()
	storeLock = make(chan (int), 1)

	Debug = _debug

	me = &Replica{
		address: "127.0.0.1",
		port:    port,
		isLive:  true,
	}

	// bind pub/sub

	pubSock, err = zmq.NewSocket(zmq.PUB)

	if err != nil {
		p_out("skt create err %v\n", err)
	}

	s := fmt.Sprintf("tcp://*:%d", PUB_PORT(port))
	err = pubSock.Bind(s)

	me.pubAddr = s

	if err != nil {
		p_out("Error binding pub/sub socket %q (%v)\n", s, err)
	}
	p_out("pub sock bound to %q\n", s)

	subSock, err = zmq.NewSocket(zmq.SUB)
	p_dieif(err != nil, "Bad SUB sock, %v\n", err)
	p_out("sub sock initialized to %q\n", subSock)

	subSock.SetSubscribe("")

	// bind req/rep
	repSock, err = zmq.NewSocket(zmq.REP)

	s = fmt.Sprintf("tcp://*:%d", REQ_PORT(port))
	if err := repSock.Bind(s); err != nil {
		p_die("Error binding req/rep socket %q (%v)\n", s, err)
	}

	p_out("request bound to %q\n", s)

	//go recvLoop()

	for {
		p_out("In repLoop\n")

		msg := recv(repSock)
		p_out("Received message\n")

		var v Value
		var err error
		var m *Msg

		switch msg.MsgType {
		case MSG_PUT:
			put(msg.Key, GetValue(msg.Value, msg.Type))
			p_out("Put %s, %v\n", msg.Key, msg.Value)
			m = &Msg{
				Key:     msg.Key,
				Value:   msg.Value,
				MsgType: MSG_PUT_RESPONSE,
				Status:  0,
			}
			break
		case MSG_GET:
			v, err = get(msg.Key)

			if err != nil {
				m = &Msg{
					Key:     msg.Key,
					MsgType: MSG_GET_RESPONSE,
					Status:  -1,
				}
			} else {
				m = &Msg{
					Key:     msg.Key,
					Value:   v.Serialize(),
					MsgType: MSG_GET_RESPONSE,
					Status:  0,
				}
			}
			break
		case MSG_DELETE:
			v, err = deleteEntry(msg.Key)

			if err != nil {
				m = &Msg{
					Key:     msg.Key,
					MsgType: MSG_GET_RESPONSE,
					Status:  -1,
				}
			} else {
				m = &Msg{
					Key:     msg.Key,
					Value:   v.Serialize(),
					MsgType: MSG_GET_RESPONSE,
					Status:  0,
				}
			}
			break
		default:
			p_out("Received unknown message type")
			break
		}
		send(repSock, m)
	}
}

// ========================================================================== //

// Wait for pubsub data (from other betareduce servers )
func recvLoop() {

	// TODO: connect to other replicas here
	p_out("In recvLoop")

	for {
		msg := recv(subSock)
		p_out("Recv msg %q\n", msg.S)
	}
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

	if e != nil {
		return nil, EKEYNF
	}

	return s, nil
}

func deleteEntry(key string) (Value, error) {
	in()
	s, e := store.Delete(key)
	out()

	if e != nil {
		return nil, EKEYNF
	}

	return s, nil
}
