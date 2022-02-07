package serialize

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

type Serializer interface {
	Marshaler
	Unmarshaler
}
