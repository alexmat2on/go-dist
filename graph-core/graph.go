package graphcore

import (
	"sync"
)

// A Graph that represents a network of Nodes. Nodes have private local data, and should
// only be able to exchange messages with their immediate neighbors in the graph 
type Graph struct {
	Nodes map[int]NodeIfc
	Edges []Edge
	Msgs chan Message
	Wg *sync.WaitGroup
}

func (g *Graph) addNode(n NodeIfc) {
	g.Nodes[n.Id()] = n
}

// Return -1 if the nodes are already neighbors (i.e., are connected by an edge) 
func (g *Graph) AddEdge(node1 NodeIfc, node2 NodeIfc) int {
	if node1.CheckNeighbors(node2) || node2.CheckNeighbors(node1) {
		return -1
	}

	node1.AddNeighbor(node2)
	g.addNode(node1)
	g.addNode(node2)

	e := Edge{n1: node1, n2: node2, comms: make(chan Message, 5), weight: 0}
	g.Edges = append(g.Edges, e)

	return 0
}

// Return the index of the edge connecting these two nodes, or -1 if there is no such edge.
func (g *Graph) FindEdge(node1 NodeIfc, node2 NodeIfc) int {
	ret := -1
	for i := 0; i < len(g.Edges); i++ {
		if g.Edges[i].Contains(node1) && g.Edges[i].Contains(node2) {
			ret = i
			break
		}
	}
	return ret
}
