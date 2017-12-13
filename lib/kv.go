package betareduce

// The KVS type represents a key value store
type KVS struct {
	store map[string]string
}

func (kv *KVS) put(key string, value string) {
	kv.store[key] = value
}

func (kv *KVS) get(key string) string {
	return kv.store[key]
}
