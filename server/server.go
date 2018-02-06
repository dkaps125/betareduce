package main

import (
	lib "betareduce/lib"
	"strconv"
	"strings"
)

var (
	store     KVS
	storeLock chan (int)
	me        *lib.Replica
)

// ========================================================================== //

func in() {
	storeLock <- 1
}

func out() {
	<-storeLock
}

// ========================================================================== //

// Init initializes a key-value store and binds to a port
func Run(port int, bootstrapAddress string) {
	store = NewKVS()
	storeLock = make(chan (int), 1)

	me = lib.CreateReplica("127.0.0.1", port)

	// bind pub/sub
	me.InitReplica()

	if bootstrapAddress != "" {
		toks := strings.Split(bootstrapAddress, ":")
		address := toks[0]
		port, _ := strconv.Atoi(toks[1])

		me.BootstrapFromReplica(address, port)
	}

	go recvLoop()

	for {
		lib.P_out("In repLoop\n")

		msg := me.RecvClient()
		lib.P_out("Received message\n")

		var v Value
		var err error
		var m *lib.Msg

		err = nil

		switch msg.MsgType {
		case lib.MSG_PUT:
			put(msg.Key, GetValue(msg.Value))
			lib.P_out("Put %s, %v\n", msg.Key, msg.Value)
			m = &lib.Msg{
				Key:     msg.Key,
				Value:   msg.Value,
				MsgType: lib.MSG_PUT_RESPONSE,
				Status:  0,
			}
			break
		case lib.MSG_GET:
			v, err = get(msg.Key)

			if err != nil {
				m = &lib.Msg{
					Key:     msg.Key,
					MsgType: lib.MSG_GET_RESPONSE,
					Status:  -1,
				}
			} else {
				m = &lib.Msg{
					Key:     msg.Key,
					Value:   v.Serialize(),
					MsgType: lib.MSG_GET_RESPONSE,
					Status:  0,
				}
			}
			break
		case lib.MSG_DELETE:
			v, err = deleteEntry(msg.Key)

			if err != nil {
				m = &lib.Msg{
					Key:     msg.Key,
					MsgType: lib.MSG_GET_RESPONSE,
					Status:  -1,
				}
			} else {
				m = &lib.Msg{
					Key:     msg.Key,
					Value:   v.Serialize(),
					MsgType: lib.MSG_GET_RESPONSE,
					Status:  0,
				}
			}
			break
		default:
			lib.P_out("Received unknown message type\n")
			err = EKEYNF
			break
		}

		if err == nil {
			me.SendClient(m)
		}
	}
}

// ========================================================================== //

// Wait for pubsub data (from other betareduce servers )
func recvLoop() {

	// TODO: connect to other replicas here
	lib.P_out("In recvLoop")

	for {
		msg := me.RecvRep()
		lib.P_out("Recv msg %q\n", msg.S)
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
