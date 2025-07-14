All code should carefully commented and be made available on GitHub with
usage instructions.
You are welcome to extend the given specification with additional
features that can improve its usability.
Whenever some input is required, its format must be unambiguously
defined.
For textual graph specifications, if possible, exploit some fragments of
the .dot language (https://graphviz.org/documentation/).

##############################################################################
#
# Dispatching network
#
Write a Google Go program that implements dispatching networks.
A dispatching network is described as a directed labelled graph such that
each node of the graph represents a sequential process and each arc a separate,
private, synchronous communication channel for transmitting strings.
Each node has a (string) label and is connected to al least one arc.
Nodes are partitioned in three categories:
- source nodes have no incoming arcs;
- target nodes have no outgoing arcs;All code should carefully commented and be made available on GitHub with
usage instructions.
You are welcome to extend the given specification with additional
features that can improve its usability.
Whenever some input is required, its format must be unambiguously
defined.
For textual graph specifications, if possible, exploit some fragments of
the .dot language (https://graphviz.org/documentation/).

##############################################################################
#
# Dispatching network
#
Write a Google Go program that implements dispatching networks.
A dispatching network is described as a directed labelled graph such that
each node of the graph represents a sequential process and each arc a separate,
private, synchronous communication channel for transmitting strings.
Each node has a (string) label and is connected to al least one arc.
Nodes are partitioned in three categories:
- source nodes have no incoming arcs;
- target nodes have no outgoing arcs;
- the remaining nodes are called internal nodes.
Each source node has also a list of messages to be sent.
Whenever an internal node receives a message from one of its incoming arcs it
appends its label to the message and forwards it along one of its outgoing arcs.
Messages received from target nodes are printed on the standard output.
The program should take in input the network description and the list of messages
to be sent by source nodes, deploy the corresponding process and run the experiment.

############################################################################## 
- the remaining nodes are called internal nodes.
Each source node has also a list of messages to be sent.
Whenever an internal node receives a message from one of its incoming arcs it
appends its label to the message and forwards it along one of its outgoing arcs.
Messages received from target nodes are printed on the standard output.
The program should take in input the network description and the list of messages
to be sent by source nodes, deploy the corresponding process and run the experiment.

############################################################################## 