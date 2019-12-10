package paralleltraversal

import (
	"fmt"
	gc "github.com/alexmat2on/go-dist/graph-core"
)

type PT_Node struct {
	gc.Node
	visited bool
}

func sendGoToAll(n *PT_Node, g *gc.Graph) {
	for i := 0; i < len(n.Neighbors()); i++ {
		other := n.Neighbors()[i]
		edgeIndex := g.FindEdge(n, other)
		if edgeIndex >= 0 {
			g.Edges[edgeIndex].Send(gc.Message{Name:"go", Sender: n.Id()})
			g.Wg.Add(1)
			go ParallelTraversalRunner(other.(*PT_Node), n.Id(), g)
		} else {
			fmt.Println("This shouldn't ever be negative...")
		}
	}
}

func (n *PT_Node) StartMsg(g *gc.Graph) {
	n.visited = true
	fmt.Println(n.Id(), ": BCT Start message received.\n")
	sendGoToAll(n, g)
}

func (n *PT_Node) GoMsg(g *gc.Graph, m *gc.Message) {
	if !n.visited {
		n.visited = true

		fmt.Println(n.Id(), ": BCT Go message received from", m.Sender, ".\n")
		sendGoToAll(n, g)
	}
}
