package broadcasttraversal

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	gc "github.com/alexmat2on/go-dist/graph-core"
)

// Generate a random graph of processing nodes
// g -- initialized graph
// numNodes -- total number of nodes the networks should contain 
// maxNeighbors -- the maximum # of neighbors per node. The final # of neighbors and the set of neighbors for each node is randomized.
// wg -- initialized waitgroup for the nodes 
func genGraph(g *gc.Graph, numNodes int, maxNeighbors int, wg *sync.WaitGroup) {
	g.Wg = wg
	g.Msgs = make(chan gc.Message, numNodes)
	g.Nodes = make(map[int]gc.NodeIfc)
	g.Edges = make([]gc.Edge, 0)

	if numNodes < 1 || maxNeighbors < 1 {
		fmt.Println("ERROR: Not enough nodes/neighbors")
	}

	if maxNeighbors > numNodes {
		fmt.Println("ERROR: Too many neighbors")
	}

	// Create a new set of nodes
	newNodes := make([]BCT_Node, numNodes)
	for i := 0; i < numNodes; i++ {
		newNodes[i] = BCT_Node{Node: gc.InitNode(i + 1), visited: false}
	}

	// Generate random edge connections between the nodes
	for i := 0; i < numNodes; i++ {
		numNeighbors := rand.Intn(maxNeighbors) + 1 // an int: 1 <= nn <= maxNeigh
		newNeighbors := rand.Perm(numNodes) // slice of ints in random order
		for j := 0; j < numNeighbors; j++ {
			if newNeighbors[j] != i {
				// Only add this neighbor if it isn't itself
				neighborIndex := newNeighbors[j]
				res := g.AddEdge(&newNodes[i], &newNodes[neighborIndex])
				if res < 0 {
					// We added an edge that already exists, try again 
					// We want to check this scenario to guarantee a connected graph
					numNeighbors = (numNeighbors + 1) % numNodes // in the case that numNeighbors is increased past the number of nodes...
				}
			} else {
				// If we happen to get the node itself, try again
				numNeighbors = (numNeighbors + 1) % numNodes
			}
		}
	}
} 

func BroadcastTraversalRunner(n *BCT_Node, senderId int, g *gc.Graph) {
	var msg gc.Message
	if senderId > 0 {
		sender := g.Nodes[senderId]
		edgeIndex := g.FindEdge(n, sender)
		if edgeIndex >= 0 {
			msg = g.Edges[edgeIndex].Receive()
		} else {
			fmt.Println("NOOOOOOO")
		}
	} else {
		msg = <-g.Msgs
	}

	switch msg.Name {
	case "start": n.StartMsg(g)
	case "go": n.GoMsg(g, &msg)
	}

	g.Wg.Done()
}

func BroadcastTraversal() {
	s := time.Now().UnixNano()
	rand.Seed(s)

	var wg sync.WaitGroup
	var g gc.Graph
	
	genGraph(&g, 10, 5, &wg)

	g.Msgs <- gc.Message{Name:"start", Sender:0}

	g.Wg.Add(1)
	go BroadcastTraversalRunner(g.Nodes[rand.Intn(len(g.Nodes)) + 1].(*BCT_Node), 0, &g)
	
	g.Wg.Wait()
}