package graph

import "errors"

// EdgeRole defines the structural role of an edge within the graph topology.
//
// The role determines how the edge participates in the graph workflow and affects
// validation and execution behavior. Each graph must have exactly one StartEdge
// and at least one EndEdge to define entry and exit points.
type EdgeRole int

var (
	// ErrSourceNodeNil indicates that the start node is nil.
	ErrSourceNodeNil = errors.New("start node cannot be nil")
	// ErrDestinationNodeNil indicates that the end node is nil.
	ErrDestinationNodeNil = errors.New("end node cannot be nil")
)

const (
	// StartEdge connects the implicit start node to the first operational node.
	//
	// Every graph execution begins by traversing the StartEdge to reach the first
	// node that will process the input. There must be exactly one StartEdge per graph.
	//
	// Created by: builders.CreateStartEdge()
	StartEdge EdgeRole = iota

	// EndEdge connects an operational node to the implicit end node.
	//
	// When execution reaches an EndEdge, the graph workflow completes and returns
	// the final state. A graph can have multiple EndEdges to support different
	// exit conditions or termination paths.
	//
	// Created by: builders.CreateEndEdge()
	EndEdge

	// IntermediateEdge connects two operational nodes in the graph.
	//
	// These edges form the main workflow paths, defining how execution flows
	// from one processing step to another. Most edges in a graph are intermediate
	// edges.
	//
	// Created by: builders.CreateEdge()
	IntermediateEdge
)

// Edge represents a directed connection between two nodes in the graph.
//
// Edges define the flow of execution through the graph workflow. When a node completes
// execution, the graph follows one or more outgoing edges (determined by the node's
// routing policy) to continue processing. Edges can carry optional labels for metadata
// and identification purposes.
//
// Edges are created using the builder functions in the builders package:
//   - builders.CreateStartEdge() for graph entry points
//   - builders.CreateEdge() for operational connections
//   - builders.CreateEndEdge() for graph exit points
type Edge[T SharedState] interface {
	// From returns the source node where this edge originates.
	//
	// The source node is the node that must complete execution before this edge
	// can be traversed. For StartEdge, this is the implicit start node.
	//
	// Returns:
	//   - The source Node of this edge.
	From() Node[T]

	// To returns the destination node where this edge terminates.
	//
	// The destination node is the node that will execute next if this edge is
	// selected by the routing policy. For EndEdge, this is the implicit end node.
	//
	// Returns:
	//   - The destination Node of this edge.
	To() Node[T]

	// LabelByKey retrieves a label value by its key from the edge's metadata.
	//
	// Labels are optional key-value pairs that can be attached to edges for
	// identification, categorization, or conditional routing logic. They are
	// provided during edge creation.
	//
	// Parameters:
	//   - key: The label key to look up.
	//
	// Returns:
	//   - The label value if the key exists.
	//   - A boolean indicating whether the key was found (true) or not (false).
	//
	// Example:
	//
	//	if label, ok := edge.LabelByKey("type"); ok {
	//	    fmt.Printf("Edge type: %s\n", label)
	//	}
	LabelByKey(key string) (string, bool)

	// Role returns the structural role of this edge in the graph.
	//
	// The role indicates whether this is a StartEdge, EndEdge, or IntermediateEdge,
	// which affects how the edge is treated during validation and execution.
	//
	// Returns:
	//   - The EdgeRole of this edge.
	Role() EdgeRole
}
