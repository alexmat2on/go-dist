package main

import (
	"fmt"
)

type Node struct {
	id int
	neighbors []int
}

func (n *Node) startMsg(g *Graph) {
	fmt.Println(n.id, ": Start message received.\n")

	g.msgs <- Message{name:"go", sender: n.id}
	g.wg.Add(1)
	nodeRunner(g.nodes[n.neighbors[0]], g)
}

func (n *Node) goMsg(g *Graph, m *Message) {
	fmt.Println(n.id, ": Go message received from", m.sender, ".\n")
}
