package betareduce

type Value interface {
	serialize() []byte
	deserialize([]byte) Value
}

// ========================================================================== //

type String struct {
	Value string
}

func (s String) serialize() []byte {
	return []byte(s.Value)
}

func (s String) deserialize(value []byte) Value {
	return String{Value: string(value)}
}
