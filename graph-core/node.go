package graphcore

// Define an interface that can generalize different kinds of nodes, different algorithms 
// require different node fields and message functions  
type NodeIfc interface {
	StartMsg(g *Graph)
	GoMsg(g *Graph, m *Message)
	Id() int
	Neighbors() []NodeIfc
	AddNeighbor(NodeIfc)
	CheckNeighbors(neigh NodeIfc) bool
}

// A "base" Node struct. Every node will need to have an identifier and neighbor set. 
type Node struct {
	id int
	neighbors []NodeIfc
}

func (n Node) Id() int {
	return n.id
}

func (n Node) Neighbors() []NodeIfc {
	return n.neighbors
}

func (n Node) CheckNeighbors(neigh NodeIfc) bool {
	for i := 0; i < len(n.neighbors); i++ {
		if n.neighbors[i] == neigh {
			return true
		}
	}
	return false
}

func (n *Node) AddNeighbor(neigh NodeIfc) {
	n.neighbors = append(n.neighbors, neigh)
}

func InitNode(id int) Node {
	return Node{id: id, neighbors: nil}
}