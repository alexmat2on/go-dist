package broadcasttraversal

import (
	"fmt"
	gc "github.com/alexmat2on/go-dist/graph-core"
)

type BCT_Node struct {
	gc.Node
	visited bool
}

func sendGoToAll(n *BCT_Node, g *gc.Graph) {
	for i := 0; i < len(n.Neighbors()); i++ {
		g.Msgs <- gc.Message{Name:"go", Sender: n.Id()}
		g.Wg.Add(1)
		go BroadcastTraversalRunner(n.Neighbors()[i].(*BCT_Node), g)
	}
}

func (n *BCT_Node) StartMsg(g *gc.Graph) {
	n.visited = true
	fmt.Println(n.Id(), ": BCT Start message received.\n")
	sendGoToAll(n, g)
}

func (n *BCT_Node) GoMsg(g *gc.Graph, m *gc.Message) {
	if !n.visited {
		n.visited = true

		fmt.Println(n.Id(), ": BCT Go message received from", m.Sender, ".\n")
		sendGoToAll(n, g)
	}
}
