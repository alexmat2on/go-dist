package graphcore

import (
	"sync"
)

// A Graph that represents a network of Nodes. Nodes have private local data, and should
// only be able to exchange messages with their immediate neighbors in the graph 
type Graph struct {
	Nodes map[int]NodeIfc
	Edges map[int][]NodeIfc
	Msgs chan Message
	Wg *sync.WaitGroup
}

func (g *Graph) AddNode(n NodeIfc) {
	g.Nodes[n.Id()] = n
}

// Return -1 if the nodes are already neighbors (i.e., are connected by an edge) 
func (g *Graph) AddEdge(node1 NodeIfc, node2 NodeIfc) int {
	if node1.CheckNeighbors(node2) || node2.CheckNeighbors(node1) {
		return -1
	}

	node1.AddNeighbor(node2)
	g.AddNode(node1)
	g.AddNode(node2)

	g.Edges[node1.Id()] = append(g.Edges[node1.Id()], node2)

	return 0
}
