package main

import (
	"encoding/binary"
)

type Value interface {
	Serialize() []byte
	Deserialize([]byte) Value
}

func GetValue(contents []byte) Value {
	var v Value
	switch contents[0] {
	case 's':
		v = String{}
		return v.Deserialize(contents[1:])
	case 'i':
		v = Int{}
		return v.Deserialize(contents[1:])
	default:
		v = ByteArray{}
		return v.Deserialize(contents[1:])
	}
}

// ========================================================================== //

type ByteArray struct {
	Value []byte
}

func (s ByteArray) Serialize() []byte {
	return append([]byte{'b'}, s.Value...)
}

func (s ByteArray) Deserialize(value []byte) Value {
	s.Value = value
	return s
}

// ========================================================================== //

type String struct {
	Value string
}

func (s String) Serialize() []byte {
	ret := []byte(s.Value)
	return append([]byte{'s'}, ret...)
}

func (s String) Deserialize(value []byte) Value {
	s.Value = string(value)
	return s
}

// ========================================================================== //

type Int struct {
	Value uint64
}

// Serialize serializes an int primitive
func (s Int) Serialize() []byte {
	ret := make([]byte, 9)
	ret[0] = 'i'
	binary.BigEndian.PutUint64(ret[1:], s.Value)
	return ret
}

func (s Int) Deserialize(value []byte) Value {
	s.Value = binary.BigEndian.Uint64(value)
	return s
}
