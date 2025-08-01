import graphviz
import sys
import os

if __name__ == "__main__":
    tmp = sys.argv[1]
    dot = graphviz.Digraph("Graph")
    
    edges = tmp[1:len(tmp)-1].split(",")
    for edge in edges:
        nodes, weight = edge.replace("[","").replace("]","").split(":")
        start, end = nodes.split("->")
        if int(weight) == 0: 
            dot.edge(start, end, "")
        else:
            dot.edge(start, end, weight)            

    dot.render('graph', format='png', view=False)
    os.remove("graph") # delete graph file that describe the graph with new labels
    print("graph created")