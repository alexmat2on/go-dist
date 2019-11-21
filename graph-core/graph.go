package graphcore

import (
	"sync"
)

type Graph struct {
	Nodes map[int]NodeInterface
	Edges map[int]NodeInterface
	Msgs chan Message
	Wg *sync.WaitGroup
}

func (g *Graph) AddNode(n NodeInterface) {
	g.Nodes[n.Id()] = n;
}

func (g *Graph) AddEdge(node1 NodeInterface, node2 NodeInterface) {
	g.AddNode(node1)
	g.AddNode(node2)
	g.Edges[node1.Id()] = node2;
}
