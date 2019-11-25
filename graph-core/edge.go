package graphcore

// Define an interface for edges
type EdgeIfc interface {
	Contains(n NodeIfc) bool
	SetWeight(w int)
	GetWeight() int
	Send(m Message)
	Receive() Message
}

// A "base" Edge struct. Every edge consists of two nodes, a communication channel 
// between them, and optionally a weight. 
type Edge struct {
	n1 NodeIfc
	n2 NodeIfc
	comms chan Message
	weight int
}

func (e *Edge) Contains(n NodeIfc) bool {
	return e.n1.Id() == n.Id() || e.n2.Id() == n.Id()
}

func (e *Edge) SetWeight(w int) {
	e.weight = w
}

func (e *Edge) GetWeight() int {
	return e.weight
}

func (e *Edge) Send(m Message) {
	e.comms <- m
}

func (e *Edge) Receive() Message {
	m := <-e.comms
	return m
}
