package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/pkg/graph"
)

// TestSharedState is a simple implementation of SharedState for testing
type TestSharedState struct {
	Value string
}

func TestCreateEdge(t *testing.T) {
	// graph.Create nodes for testing
	fromNode := graph.CreateStartNode[TestSharedState]()
	toNode := graph.CreateEndNode[TestSharedState]()

	// Test graph.CreateEdge
	edge := graph.CreateEdge(fromNode, toNode)

	if edge == nil {
		t.Fatal("graph.CreateEdge returned nil")
	}

	if edge.To() != toNode {
		t.Error("Edge.To() did not return the expected node")
	}
}

func TestCreateStartEdge(t *testing.T) {
	// graph.Create a target node
	toNode := graph.CreateEndNode[TestSharedState]()

	// Test graph.CreateStartEdge
	startEdge := graph.CreateStartEdge(toNode)

	if startEdge == nil {
		t.Fatal("graph.CreateStartEdge returned nil")
	}

	if startEdge.To() != toNode {
		t.Error("StartEdge.To() did not return the expected node")
	}
}

func TestCreateEndEdge(t *testing.T) {
	// graph.Create a source node
	fromNode := graph.CreateStartNode[TestSharedState]()

	// Test graph.CreateEndEdge
	endEdge := graph.CreateEndEdge(fromNode)

	if endEdge == nil {
		t.Fatal("graph.CreateEndEdge returned nil")
	}

	// The To() method should return an EndNode
	toNode := endEdge.To()
	if _, ok := toNode.(*graph.EndNode[TestSharedState]); !ok {
		t.Error("EndEdge.To() did not return an EndNode")
	}
}

func TestStartEdgeTo(t *testing.T) {
	toNode := graph.CreateEndNode[TestSharedState]()
	startEdge := graph.CreateStartEdge(toNode)

	result := startEdge.To()
	if result != toNode {
		t.Error("StartEdge.To() returned wrong node")
	}
}

func TestEndEdgeTo(t *testing.T) {
	fromNode := graph.CreateStartNode[TestSharedState]()
	endEdge := graph.CreateEndEdge(fromNode)

	result := endEdge.To()
	if _, ok := result.(*graph.EndNode[TestSharedState]); !ok {
		t.Error("EndEdge.To() should return an EndNode")
	}
}

func TestEdgeImplTo(t *testing.T) {
	fromNode := graph.CreateStartNode[TestSharedState]()
	toNode := graph.CreateEndNode[TestSharedState]()

	edge := graph.CreateEdge(fromNode, toNode)

	result := edge.To()
	if result != toNode {
		t.Error("edgeImpl.To() returned wrong node")
	}
}

func TestEdgeInterfaceCompliance(t *testing.T) {
	// Test that all edge types implement the Edge interface
	toNode := graph.CreateEndNode[TestSharedState]()
	fromNode := graph.CreateStartNode[TestSharedState]()

	var edge graph.Edge[TestSharedState]

	// Test StartEdge
	startEdge := graph.CreateStartEdge(toNode)
	edge = startEdge
	if edge.To() != toNode {
		t.Error("StartEdge does not properly implement Edge interface")
	}

	// Test EndEdge
	endEdge := graph.CreateEndEdge(fromNode)
	edge = endEdge
	if edge.To() == nil {
		t.Error("EndEdge does not properly implement Edge interface")
	}

	// Test edgeImpl
	edgeImpl := graph.CreateEdge(fromNode, toNode)
	edge = edgeImpl
	if edge.To() != toNode {
		t.Error("edgeImpl does not properly implement Edge interface")
	}
}
