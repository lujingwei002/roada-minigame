package cli

type (
	LifetimeHandler func(*Session)

	Lifetime struct {
		closeHandlers []LifetimeHandler
		openHandlers  []LifetimeHandler
	}
)

func newLifetime() *Lifetime {
	return &Lifetime{
		closeHandlers: make([]LifetimeHandler, 0),
		openHandlers:  make([]LifetimeHandler, 0),
	}
}

func (self *Lifetime) onClose(h LifetimeHandler) {
	self.closeHandlers = append(self.closeHandlers, h)
}

func (self *Lifetime) close(s *Session) {
	if len(self.closeHandlers) < 1 {
		return
	}
	for _, h := range self.closeHandlers {
		h(s)
	}
}

func (self *Lifetime) onOpen(h LifetimeHandler) {
	self.openHandlers = append(self.openHandlers, h)
}

func (self *Lifetime) open(s *Session) {
	if len(self.openHandlers) < 1 {
		return
	}
	for _, h := range self.openHandlers {
		h(s)
	}
}
