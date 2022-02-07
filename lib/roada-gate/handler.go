package gat

type HandlerInterface interface {
	ServeMessage(r *Request)
	OnSessionOpen(s *Session)
	OnSessionClose(s *Session)
}

type MiddleWareInterface interface {
	ServeMessage(r *Request)
	OnSessionOpen(s *Session)
	OnSessionClose(s *Session)
}
