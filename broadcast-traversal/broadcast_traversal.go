package broadcasttraversal

import (
	"sync"
	gc "github.com/alexmat2on/go-dist/graph-core"
)

func BroadcastTraversalRunner(n gc.NodeInterface, g *gc.Graph) {
	m := <-g.Msgs
	switch m.Name {
	case "start": n.StartMsg(g)
	case "go": n.GoMsg(g, &m)
	}

	g.Wg.Done()
}

func BroadcastTraversal() {
	var wg sync.WaitGroup
	var g gc.Graph
	g.Wg = &wg
	g.Msgs = make(chan gc.Message, 4)
	g.Nodes = make(map[int]gc.NodeInterface)
	g.Edges = make(map[int]gc.NodeInterface)

	g.Msgs <- gc.Message{Name:"start", Sender:0}

	n1 := BCT_Node{Node: gc.InitNode(1, []int{2, 4}), visited: false}
	n2 := BCT_Node{Node: gc.InitNode(2, []int{1, 3}), visited: false}
	n3 := BCT_Node{Node: gc.InitNode(3, []int{2}), visited: false}
	n4 := BCT_Node{Node: gc.InitNode(4, []int{1}), visited: false}

	g.AddEdge(&n1, &n2);
	g.AddEdge(&n1, &n4);
	g.AddEdge(&n2, &n1);
	g.AddEdge(&n2, &n3);
	g.AddEdge(&n3, &n2);
	g.AddEdge(&n4, &n1);

	g.Wg.Add(1)
	go BroadcastTraversalRunner(&n2, &g)
	
	g.Wg.Wait()
}