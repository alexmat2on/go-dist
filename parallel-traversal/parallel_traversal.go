package paralleltraversal

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
	gc "github.com/alexmat2on/go-dist/graph-core"
)

// Given a slice and a value, remove the value from that slice. 
func delElem(s []int, v int) []int {
	ret := make([]int, 0)
	
	for i := 0; i < len(s); i++ {
		if s[i] != v {
			ret = append(ret, s[i])
		}
	}

	return ret
}

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
	newNodes := make([]PT_Node, numNodes)
	for i := 0; i < numNodes; i++ {
		newNodes[i] = PT_Node{Node: gc.InitNode(i + 1), visited: false}

		// Give each node a randomized degree
		degree := rand.Intn(maxNeighbors)
		if degree < 2 {
			degree = 2 // minimum node degree of 2
		}

		newNodes[i].SetDegree(degree)
	}

	// Generate random edge connections between the nodes.
	// Idea: keep two sets -- in_graph, an array of newNode indices that are in the connected graph G
	// and out_graph, an array of newNode indices that are not included yet. 
	//
	// For every new Node, it's first neighbor should be picked randomly from in_graph. 
	// If it has space for more neighbors, then pick randomly from the out_graph and move 
	// the selected neighbor into the in_graph
	//
	// Repeat until the out_graph is empty. 

	// Get a random sequence of nodes to add as neighbors. 
	out_graph := rand.Perm(numNodes)
	original_set := out_graph

	// Put the first random node into in_graph
	in_graph := []int{out_graph[0]}
	newNodes[out_graph[0]].SetDegree(1)
	out_graph = delElem(out_graph, out_graph[0]) // delete by value

	// start i at 1 since we handled the first node
	for i := 1; i < numNodes; i++ { 
		node := &newNodes[original_set[i]]

		// First neighbor should be random from in_graph
		neighIdx := in_graph[rand.Intn(len(in_graph))]

		fmt.Println("Picking first neighbor...")
		// Make sure the chosen node has room for more neighbors
		for {
			if len(out_graph) == 0 {
				break
			}
			if len(newNodes[neighIdx].Neighbors()) < newNodes[neighIdx].Degree() {
				g.AddEdge(node, &newNodes[neighIdx])
				break
			} else {
				// Sample a new node
				neighIdx = in_graph[rand.Intn(len(in_graph))]
			}	
		}

		// Next neighbors should be picked randomly from out_graph, until degree is filled
		fmt.Println("Picking next neighbors...")
		j := 0

		degree := node.Degree() - 1 // Leave one slot available for future nodes to hook into. 
		for degree > 0 {
			if len(out_graph) == 0 {
				break
			}

			neighIdx = out_graph[j]

			g.AddEdge(node, &newNodes[neighIdx])
			out_graph = delElem(out_graph, neighIdx)
			in_graph = append(in_graph, neighIdx)
			degree = degree - 1

			if len(out_graph) == 0 {
				// Why do we need this twice? Idk :( 
				break
			}

			j = (j + 1) % len(out_graph)
		}
	}
} 

func ParallelTraversalRunner(n *PT_Node, senderId int, g *gc.Graph) {
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

func ParallelTraversal() {
	s := time.Now().UnixNano()
	rand.Seed(s)

	var wg sync.WaitGroup
	var g gc.Graph
	
	genGraph(&g, 10, 5, &wg)

	fmt.Println("GRAPH =============")
	fmt.Println(g.Nodes)
	g.PrintEdges()
	fmt.Println("====================")

	g.Msgs <- gc.Message{Name:"start", Sender:0}

	g.Wg.Add(1)
	initiator := rand.Intn(len(g.Nodes)) + 1
	go ParallelTraversalRunner(g.Nodes[initiator].(*PT_Node), 0, &g)
	
	g.Wg.Wait()
}