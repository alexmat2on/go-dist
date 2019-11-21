package main

import (
	"sync"
)

type Graph struct {
	nodes map[int]*Node
	edges map[int]*Node
	msgs chan Message
	wg *sync.WaitGroup
}

func (g *Graph) AddNode(n *Node) {
	g.nodes[n.id] = n;
}

func (g *Graph) AddEdge(node1 *Node, node2 *Node) {
	g.AddNode(node1)
	g.AddNode(node2)
	g.edges[node1.id] = node2;
}
