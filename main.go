package main

import (
	"sync"
)

func nodeRunner(n *Node, g *Graph) {
	m := <-g.msgs
	switch m.name {
	case "start": n.startMsg(g)
	case "go": n.goMsg(g, &m)
	}

	g.wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var g Graph
	g.wg = &wg
	g.msgs = make(chan Message, 4)
	g.nodes = make(map[int]*Node)
	g.edges = make(map[int]*Node)

	g.msgs <- Message{name:"start", sender:0}

	n1 := Node{id: 1, neighbors: []int{2, 4}}
	n2 := Node{id: 2, neighbors: []int{1, 3}}
	n3 := Node{id: 3, neighbors: []int{2}}
	n4 := Node{id: 4, neighbors: []int{1}}

	g.AddEdge(&n1, &n2);
	g.AddEdge(&n1, &n4);
	g.AddEdge(&n2, &n1);
	g.AddEdge(&n2, &n3);
	g.AddEdge(&n3, &n2);
	g.AddEdge(&n4, &n1);

	g.wg.Add(1)
	go nodeRunner(&n3, &g)
	
	g.wg.Wait()
}