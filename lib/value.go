package betareduce

import "strconv"

type Value interface {
	Serialize() []byte
	Deserialize([]byte) Value
}

func GetValue(contents []byte, kind string) Value {
	var v Value
	if kind == "String" {
		v = String{}
		return v.Deserialize(contents)
	} else if kind == "Int" {
		v = Int{}
		return v.Deserialize(contents)
	}

	return nil
}

// ========================================================================== //

type String struct {
	Value string
}

func (s String) Serialize() []byte {
	return []byte(s.Value)
}

func (s String) Deserialize(value []byte) Value {
	s.Value = string(value)
	return s
}

// ========================================================================== //

type Int struct {
	Value int
}

func (s Int) Serialize() []byte {
	return []byte(strconv.Itoa(s.Value))
}

func (s Int) Deserialize(value []byte) Value {
	s.Value, _ = strconv.Atoi(string(value))
	return s
}
