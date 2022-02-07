package gat

type Request struct {
	Session *Session
	Route   string
	Payload []byte
	mid     uint64
}

func (r *Request) Response(v interface{}) error {
	return r.Session.Response(r.mid, v)
}
