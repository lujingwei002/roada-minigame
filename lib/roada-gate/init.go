package gat

var defaultGate *Gate

func Default() *Gate {
	return defaultGate
}

func init() {
	defaultGate = NewGate()
}
