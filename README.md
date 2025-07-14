# Dispatching Network by Andrea Mussari

## Prerequisites
`graphviz` library for python

## Implementation
Constants:
- Undefined = 0
- Source = 1
- Internal = 2
- Target = 3

Structs:
- Node rappresent a node of the graph.
    - Label: the node's name
    - Type: the node's type (constants)
    - Incoming: slice of channel that represents the edges incoming
    - Outgoing: slice of channel that represents the edges outgoing
    - Messages: slice of strings only used by source nodes

- Edge rappresent a edge of the graph.
    - Label: `A->B`
    - Ch: channel shared between two nodes
    - Traffic: number of messages that pass through that edge

- NetworkGraph rappresent the graph.
    - Nodes: slice of Nodes
    - Edges: slice of Edges
    - Messages: Slice of string that contains all messages
    - Wg: WaitGroup used to is used to wait for all the goroutines launched.

Files:
1. `Exercise.md` contains the text of exercise.
2. `graph.dot` contains the description of the graph.
3. `messages.txt` contains all messages of the graph.
    - each line represents a message.
4. `script.py` script in python used to create `graph.png`.
5. `graph.png` represents the graph with weighted edges