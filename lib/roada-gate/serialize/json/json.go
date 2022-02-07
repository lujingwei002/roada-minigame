package json

import (
	"encoding/json"
)

type Serializer struct{}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
