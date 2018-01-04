package betareduce

type Value interface {
	serialize() []byte
	deserialize([]byte) Value
}

// ========================================================================== //

type String struct {
	value string
}

func (s String) serialize() []byte {
	return []byte(s.value)
}

func (s String) deserialize(value []byte) Value {
	return String{value: string(value)}
}
