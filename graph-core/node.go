package graphcore

type NodeInterface interface {
	StartMsg(g *Graph)
	GoMsg(g *Graph, m *Message)
	Id() int
	Neighbors() []int
}

type Node struct {
	id int
	neighbors []int
}

func (n *Node) Id() int {
	return n.id
}

func (n *Node) Neighbors() []int {
	return n.neighbors
}

func InitNode(id int, neighbors []int) Node {
	return Node{id: id, neighbors: neighbors}
}