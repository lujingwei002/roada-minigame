package cli

type Request struct {
	Session *Session
	Route   string
	mid     uint64
	Payload []byte
}
