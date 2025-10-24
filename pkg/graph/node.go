package graph

import "errors"

var (
	// ErrReservedNodeName indicates that the node name is reserved and cannot be used.
	ErrReservedNodeName = errors.New("node name is reserved and cannot be used")
)

// NodeRole represents the structural role of a node within the graph topology.
//
// The role determines how the node participates in the graph workflow and affects
// validation and execution behavior. StartNode and EndNode are implicit nodes
// managed by the runtime, while IntermediateNode represents user-created operational nodes.
type NodeRole int

const (
	// StartNode is the implicit entry point of the graph.
	//
	// This is an internal node automatically managed by the runtime. It has no
	// processing function and serves only as the source for the StartEdge.
	// Users do not create StartNode directly; it's created internally when
	// using builders.CreateStartEdge().
	StartNode NodeRole = iota

	// IntermediateNode is an operational node that processes state.
	//
	// These are the main processing units in the graph, created by users to
	// implement workflow logic. They execute NodeFn functions and use routing
	// policies to determine the next execution path.
	//
	// Created by: builders.CreateNode() or builders.CreateNodeWithRoutingPolicy()
	IntermediateNode

	// EndNode is the implicit exit point of the graph.
	//
	// This is an internal node automatically managed by the runtime. It has no
	// processing function and serves only as the destination for EndEdges.
	// Users do not create EndNode directly; it's created internally when
	// using builders.CreateEndEdge().
	EndNode
)

// Node represents a processing unit or decision point in the graph workflow.
//
// Nodes are the fundamental building blocks of graph execution. Each node can:
//   - Execute a processing function (NodeFn) to transform state
//   - Apply a routing policy to select the next edge(s) to follow
//   - Emit partial state updates during long-running operations
//
// Nodes are connected by edges to form a directed graph that defines the workflow.
// When a node completes execution, its routing policy determines which outgoing
// edge(s) to traverse next.
//
// Nodes are created using builder functions in the builders package:
//   - builders.CreateNode() for standard nodes with default routing
//   - builders.CreateNodeWithRoutingPolicy() for nodes with custom routing
//   - builders.CreateRouter() for routing-only nodes without processing logic
type Node[T SharedState] interface {
	// Accept executes this node's processing logic with the given input and runtime context.
	//
	// This method is called by the runtime when execution reaches this node. It
	// executes the node's NodeFn (if present), updates the state, and coordinates
	// with the StateObserver to track execution progress.
	//
	// This is an internal method used by the runtime; users typically do not call
	// this directly. Instead, nodes are executed automatically by Runtime.Invoke().
	//
	// Parameters:
	//   - userInput: The original input provided to Runtime.Invoke().
	//   - runtime: The StateObserver that tracks state changes and execution flow.
	Accept(userInput T, runtime StateObserver[T])

	// Name returns the unique identifier for this node.
	//
	// The name is used for identification in state monitoring entries, debugging,
	// and graph visualization. It must be unique within the graph.
	//
	// Returns:
	//   - The node's name as specified during creation.
	Name() string

	// RoutePolicy returns the routing policy that determines edge selection.
	//
	// The routing policy is invoked after this node completes execution to decide
	// which outgoing edge(s) to follow. Different policies enable different workflow
	// patterns like conditional branching, parallel execution, or loops.
	//
	// Returns:
	//   - The RoutePolicy associated with this node.
	RoutePolicy() RoutePolicy[T]

	// Role returns the structural role of this node in the graph.
	//
	// The role indicates whether this is a StartNode, IntermediateNode, or EndNode,
	// which affects how the node is treated during validation and execution.
	//
	// Returns:
	//   - The NodeRole of this node.
	Role() NodeRole
}
