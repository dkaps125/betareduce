package store

type Value interface {
	serialize() string
	deserialize(string) Value
}

// ========================================================================== //

type String struct {
	value string
}

func (s String) serialize() string {
	return s.value
}

func (s String) deserialize(value string) Value {
	return String{value: value}
}
