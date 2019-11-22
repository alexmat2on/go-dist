package broadcasttraversal

import (
	"sync"
	gc "github.com/alexmat2on/go-dist/graph-core"
	"math/rand"
	"time"
	"fmt"
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
	g.Edges = make(map[int][]gc.NodeIfc)

	if numNodes < 1 || maxNeighbors < 1 {
		fmt.Println("ERROR: Not enough nodes/neighbors")
	}

	if maxNeighbors > numNodes {
		fmt.Println("ERROR: Too many neighbors")
	}

	newNodes := make([]BCT_Node, numNodes)
	for i := 0; i < numNodes; i++ {
		newNodes[i] = BCT_Node{Node: gc.InitNode(i + 1), visited: false}
	}

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
					numNeighbors = (numNeighbors + 1) % numNodes // in the case that numNeighbors is increased past the number of nodes...
				}
			} else {
				// If we happen to get the node itself, try again
				numNeighbors = (numNeighbors + 1) % numNodes
			}
		}
	}
} 

func BroadcastTraversalRunner(n *BCT_Node, g *gc.Graph) {
	m := <-g.Msgs
	switch m.Name {
	case "start": n.StartMsg(g)
	case "go": n.GoMsg(g, &m)
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
	go BroadcastTraversalRunner(g.Nodes[1].(*BCT_Node), &g)
	
	g.Wg.Wait()
}