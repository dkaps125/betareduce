package store

var (
	store     KVS
	storeLock chan (int)
)

// ========================================================================== //

func in() {
	storeLock <- 1
}

func out() {
	<-storeLock
}

// ========================================================================== //

// Init initializes a key-value store
func Init(port int) {
	store = NewKVS()
	storeLock = make(chan (int), 1)
}

func Run() {

}

// ========================================================================== //
//Key value function wrappers

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
