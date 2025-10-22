package builders

import (
	i "github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// CreateEdge creates a new intermediate edge connecting two nodes in a graph.
//
// This function constructs an edge from a source node to a destination node,
// establishing a directed connection in the graph workflow. Optional labels can
// be provided as key-value pairs to annotate the edge with metadata.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - from: The source node where the edge originates.
//   - to: The destination node where the edge terminates.
//   - labels: Optional maps of string key-value pairs for edge metadata/annotations.
//
// Returns:
//   - A new Edge instance connecting the specified nodes.
//
// Example:
//
//	node1, _ := CreateNode[MyState]("node1", myFunction)
//	node2, _ := CreateNode[MyState]("node2", anotherFunction)
//	edge := CreateEdge(node1, node2, map[string]string{"type": "conditional"})
func CreateEdge[T g.SharedState](from, to g.Node[T], labels ...map[string]string) g.Edge[T] {
	return i.EdgeImplFactory(from, to, g.IntermediateEdge, labels...)
}

// CreateStartEdge creates a new edge from the implicit start node to a specified node.
//
// This function is used to define the entry point of a graph workflow by connecting
// the internal start node to the first operational node. Every graph execution begins
// at the start node, and this edge determines which node receives the initial state.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - to: The first operational node in the graph that will receive the initial state.
//
// Returns:
//   - A new StartEdge instance connecting the implicit start node to the specified node.
//
// Example:
//
//	firstNode, _ := CreateNode[MyState]("first", myFunction)
//	startEdge, _ := CreateStartEdge(firstNode)
func CreateStartEdge[T g.SharedState](to g.Node[T]) g.Edge[T] {
	startNode, _ := createStartNode[T]()
	return i.EdgeImplFactory(startNode, to, g.StartEdge)
}

// CreateEndEdge creates a new edge from a specified node to the implicit end node.
//
// This function is used to define an exit point of a graph workflow by connecting
// an operational node to the internal end node. When execution reaches an end edge,
// the graph workflow completes and returns the final state.
//
// Type Parameters:
//   - T: The SharedState type that will be passed through the graph execution.
//
// Parameters:
//   - from: The operational node from which the graph workflow will terminate.
//   - labels: Optional maps of string key-value pairs for edge metadata/annotations.
//
// Returns:
//   - A new EndEdge instance connecting the specified node to the implicit end node.
//
// Example:
//
//	lastNode, _ := CreateNode[MyState]("last", myFunction)
//	endEdge, _ := CreateEndEdge(lastNode)
func CreateEndEdge[T g.SharedState](from g.Node[T], labels ...map[string]string) g.Edge[T] {
	endNode, _ := createEndNode[T]()
	return i.EdgeImplFactory(from, endNode, g.EndEdge, labels...)
}
