package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

// Types of a Node
const (
	Undefined uint = iota
	Source
	Internal
	Target
)

// Edge represents a edge betwenn two nodes of the network graph
type Edge struct {
	Label   string
	Ch      chan string
	Traffic int
}

// Node represents a node in the network graph
type Node struct {
	Label    string
	Type     uint
	Incoming []chan string
	Outgoing []chan string
	Messages []string // Only used for source nodes
}

// NetworkGraph represents the entire network
type NetworkGraph struct {
	Nodes    []*Node
	Edges    []*Edge
	Messages []string
	Wg       sync.WaitGroup
}

// function to avoid more than 1 occurence of a Node in G.Nodes
func (g *NetworkGraph) hasNode(label string) bool {
	for _, node := range g.Nodes {
		if node.Label == label {
			return true
		}
	}
	return false
}

// function to get the Edge which Edge.ch is ch
func (G *NetworkGraph) getEdge(ch chan string) *Edge {
	for _, Edge := range G.Edges {
		if Edge.Ch == ch {
			return Edge
		}
	}
	return nil
}

// function used to send messages from source node
func (G *NetworkGraph) RunSourceNode(node *Node) {
	defer G.Wg.Done()

	och := 0 // index of current outgoing channel
	im := 0  // index of current message

	// loop to split all messages using all node's outgoing channels
	for {

		if im >= len(node.Messages) {
			break
		}

		if och >= len(node.Outgoing) {
			och = 0
		}

		node.Outgoing[och] <- node.Messages[im] // Blocking send to ensure delivery
		//fmt.Printf("Source %s sent: %s\n", node.Label, message)

		och++
		im++
	}

	// Close outgoing channels after sending all messages
	for _, ch := range node.Outgoing {
		close(ch)
	}
}

// function used to forward messages from source node
func (G *NetworkGraph) RunInternalNode(node *Node) {
	defer G.Wg.Done()

	var wg sync.WaitGroup
	var mu sync.Mutex // mutex used to synchronize selection of current outgoing channel
	och := 0          // index of current outgoing channel

	// Start a goroutine for each incoming channel
	for _, inCh := range node.Incoming {
		wg.Add(1)

		go func(ch chan string) {
			defer wg.Done()

			edge := G.getEdge(ch)

			for message := range ch {

				edge.Traffic++

				mu.Lock()

				if och >= len(node.Outgoing) {
					och = 0
				}

				outCh := node.Outgoing[och]
				och++

				mu.Unlock()

				outCh <- message
			}
		}(inCh)
	}

	// Wait for all processes to complete
	wg.Wait()

	// Close outgoing channels
	for _, ch := range node.Outgoing {
		close(ch)
	}
}

// function used to receive all messages from other nodes
func (G *NetworkGraph) RunTargetNode(node *Node) {
	defer G.Wg.Done()

	var wg sync.WaitGroup
	var mu sync.Mutex // mutex used to synchronize the print of messages received

	// Start a goroutine for each incoming channel
	for _, inCh := range node.Incoming {
		wg.Add(1)

		go func(ch chan string) {
			defer wg.Done()

			edge := G.getEdge(ch)

			for message := range ch {

				edge.Traffic++

				mu.Lock()

				fmt.Printf("TARGET OUTPUT: %s\n", message)

				mu.Unlock()
			}
		}(inCh)
	}

	// Wait for all processes to complete
	wg.Wait()
}

// function used to start the dispatching network
func (G *NetworkGraph) Run() {
	fmt.Println("Starting dispatching network...")

	// Start all node processes
	for _, node := range G.Nodes {
		G.Wg.Add(1)
		switch node.Type {
		case Source:
			go G.RunSourceNode(node)
		case Internal:
			go G.RunInternalNode(node)
		case Target:
			go G.RunTargetNode(node)
		}
	}

	// Wait for all processes to complete
	G.Wg.Wait()
	fmt.Println("Dispatching network completed.")
}

// function used to extract all messages from messages.txt
func GetMessages() []string {
	messages := []string{}
	fm, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer fm.Close()

	Mscanner := bufio.NewScanner(fm)

	for Mscanner.Scan() {
		line := Mscanner.Text()
		messages = append(messages, line)
	}

	return messages
}

// function used to create a new network graph
func CreateNetworkGraph() *NetworkGraph {

	G := NetworkGraph{
		Nodes:    []*Node{},
		Edges:    []*Edge{},
		Messages: GetMessages(),
		Wg:       sync.WaitGroup{},
	}

	fg, err := os.Open("graph.dot")

	if err != nil {
		log.Fatal(err)
	}

	defer fg.Close()

	Gscanner := bufio.NewScanner(fg)

	for Gscanner.Scan() {
		line := Gscanner.Text()
		if strings.Contains(line, "->") {
			edge := strings.Split(line, "->")

			edge[0] = strings.ReplaceAll(edge[0], " ", "")
			edge[1] = strings.ReplaceAll(edge[1], " ", "")

			from := Node{Label: edge[0], Type: Undefined, Incoming: make([]chan string, 0), Outgoing: make([]chan string, 0), Messages: make([]string, 0)}
			to := Node{Label: edge[1], Type: Undefined, Incoming: make([]chan string, 0), Outgoing: make([]chan string, 0), Messages: make([]string, 0)}

			G.Edges = append(G.Edges, &Edge{Label: edge[0] + "->" + edge[1], Traffic: 0})

			if !G.hasNode(from.Label) {
				G.Nodes = append(G.Nodes, &from)
			}

			if !G.hasNode(to.Label) {
				G.Nodes = append(G.Nodes, &to)
			}
		}
	}

	// loop used to fill outgoing and incoming slices
	for _, Edge := range G.Edges {
		edge := strings.Split(Edge.Label, "->")

		from := slices.IndexFunc(G.Nodes, func(n *Node) bool {
			return n.Label == edge[0]
		})

		to := slices.IndexFunc(G.Nodes, func(n *Node) bool {
			return n.Label == edge[1]
		})

		ch := make(chan string)
		Edge.Ch = ch

		// Add as outgoing channel to source node
		G.Nodes[from].Outgoing = append(G.Nodes[from].Outgoing, ch)

		// Add as incoming channel to target node
		G.Nodes[to].Incoming = append(G.Nodes[to].Incoming, ch)
	}

	// determine the type of each node
	Sources := []*Node{}
	for _, node := range G.Nodes {

		if len(node.Incoming) == 0 {
			node.Type = Source
			Sources = append(Sources, node)
		} else if len(node.Outgoing) == 0 {
			node.Type = Target
		} else {
			node.Type = Internal
		}
	}

	mxnode := int(len(G.Messages) / len(Sources)) // number of messages for source node

	// loop used to split homogeneously between all source nodes
	i := 0
	j := 0
	for {

		if i >= len(Sources) {
			i = 0
		}

		node := Sources[i]

		if len(G.Messages)-j < mxnode {
			node.Messages = append(node.Messages, G.Messages[j:]...)
			break

		} else {
			node.Messages = append(node.Messages, G.Messages[j:j+mxnode]...)
			j += mxnode
		}

		i++
	}

	return &G

}

// function used to create 'graph.png' using 'script.py'
func (G *NetworkGraph) PrintGraph() {
	tmp := "["

	for i, Edge := range G.Edges {
		if i < len(G.Edges)-1 {
			tmp += fmt.Sprintf("[%s:%x],", Edge.Label, Edge.Traffic)
		} else {
			tmp += fmt.Sprintf("[%s:%x]", Edge.Label, Edge.Traffic)
		}
	}

	tmp += "]"

	cmd := exec.Command("python3", "script.py", tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python code:", err)
		return
	}
	fmt.Print(string(out))
}

func main() {
	fmt.Println("Dispatching Network Implementation")
	fmt.Println("==================================")
	fmt.Printf("Source:%x, Internal:%x, Target:%x\n\n", Source, Internal, Target)

	G := CreateNetworkGraph()

	fmt.Println("List of Nodes:")
	for _, node := range G.Nodes {
		if node.Type == Source {
			fmt.Printf("Node %s: type=%x messages=%x\n", node.Label, node.Type, len(node.Messages))
		} else {
			fmt.Printf("Node %s: type=%x\n", node.Label, node.Type)
		}
	}

	fmt.Println("\nNetwork topology:")
	for _, edge := range G.Edges {
		fmt.Println(edge.Label)
	}
	fmt.Println("==================================")

	G.Run()

	G.PrintGraph()
}
