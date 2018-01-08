package main

import (
	"time"

	lib "betareduce/lib"
)

// ========================================================================== //

type keynotfound struct {
	s string
}

func (e keynotfound) Error() string {
	return e.s
}

// EKEYNF is an error thrown when a key is not found
var EKEYNF = keynotfound{s: "Key not found"}

// ========================================================================== //

type entry struct {
	atime time.Time
	mtime time.Time
	value Value
}

// ========================================================================== //

// The KVS type represents a key value store
type KVS struct {
	store map[string]entry
}

// NewKVS returns an instance of a key value store
func NewKVS() KVS {
	return KVS{store: make(map[string]entry)}
}

func (kv *KVS) Put(key string, value Value) {
	lib.P_out("%v\n", kv.store)
	t := time.Now()

	if _, ok := kv.store[key]; ok {
		kv.store[key] = entry{
			atime: kv.store[key].atime,
			mtime: t,
			value: value,
		}
	} else {
		kv.store[key] = entry{
			atime: t,
			mtime: t,
			value: value,
		}
	}
}

func (kv *KVS) Get(key string) (Value, error) {
	if _, ok := kv.store[key]; ok {
		return kv.store[key].value, nil
	}

	return nil, EKEYNF
}

func (kv *KVS) Delete(key string) (Value, error) {
	if _, ok := kv.store[key]; ok {
		v := kv.store[key].value
		delete(kv.store, key)
		return v, nil
	}

	return nil, EKEYNF
}
