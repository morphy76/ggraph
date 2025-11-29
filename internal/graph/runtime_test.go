package graph

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// RuntimeTestState is a simple state type for testing
type RuntimeTestState struct {
	Value   string
	Counter int
	Data    map[string]interface{}
}

// Mock node implementation for testing
type mockRuntimeNode struct {
	name      string
	role      g.NodeRole
	fn        g.NodeFn[RuntimeTestState]
	policy    g.RoutePolicy[RuntimeTestState]
	callCount int
	mu        sync.Mutex
	mailbox   chan RuntimeTestState
}

func newMockRuntimeNode(name string, role g.NodeRole, fn g.NodeFn[RuntimeTestState], policy g.RoutePolicy[RuntimeTestState]) *mockRuntimeNode {
	return &mockRuntimeNode{
		name:    name,
		role:    role,
		fn:      fn,
		policy:  policy,
		mailbox: make(chan RuntimeTestState, 10),
	}
}

func (n *mockRuntimeNode) Accept(userInput RuntimeTestState, runtime g.StateObserver[RuntimeTestState], config g.InvokeConfig) {
	go func() {
		n.mu.Lock()
		n.callCount++
		n.mu.Unlock()

		// Wait for message in mailbox
		asyncInput := <-n.mailbox

		currentState := runtime.CurrentState(config.ThreadID)

		if n.fn != nil {
			newState, err := n.fn(asyncInput, currentState, func(partial RuntimeTestState) {
				runtime.NotifyStateChange(n, config, userInput, partial, Replacer[RuntimeTestState], nil, true)
			})
			runtime.NotifyStateChange(n, config, userInput, newState, Replacer[RuntimeTestState], err, false)
		} else {
			// Router node - just pass through current state
			runtime.NotifyStateChange(n, config, userInput, currentState, Replacer[RuntimeTestState], nil, false)
		}
	}()

	// Send to mailbox
	n.mailbox <- userInput
}

func (n *mockRuntimeNode) Name() string {
	return n.name
}

func (n *mockRuntimeNode) RoutePolicy() g.RoutePolicy[RuntimeTestState] {
	return n.policy
}

func (n *mockRuntimeNode) Role() g.NodeRole {
	return n.role
}

func (n *mockRuntimeNode) GetCallCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.callCount
}

// Mock edge implementation for testing
type mockRuntimeEdge struct {
	from   g.Node[RuntimeTestState]
	to     g.Node[RuntimeTestState]
	role   g.EdgeRole
	labels map[string]string
}

func (e *mockRuntimeEdge) From() g.Node[RuntimeTestState] {
	return e.from
}

func (e *mockRuntimeEdge) To() g.Node[RuntimeTestState] {
	return e.to
}

func (e *mockRuntimeEdge) Role() g.EdgeRole {
	return e.role
}

func (e *mockRuntimeEdge) LabelByKey(key string) (string, bool) {
	val, ok := e.labels[key]
	return val, ok
}

// TestRuntimeFactory_BasicCreation tests creating a runtime with valid start edge
func TestRuntimeFactory_BasicCreation(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	initialState := RuntimeTestState{Value: "initial", Counter: 0}

	runtime, err := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: initialState})
	if err != nil {
		t.Fatalf("RuntimeFactory() failed: %v", err)
	}
	defer runtime.Shutdown()

	if runtime == nil {
		t.Fatal("RuntimeFactory() returned nil runtime")
	}

	if runtime.StartEdge() != startEdge {
		t.Error("StartEdge() did not return the provided start edge")
	}
}

// TestRuntimeFactory_NilStartEdge tests that creating runtime with nil start edge fails
func TestRuntimeFactory_NilStartEdge(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	initialState := RuntimeTestState{Value: "initial"}

	runtime, err := RuntimeFactory[RuntimeTestState](nil, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: initialState})
	if err == nil {
		t.Fatal("Expected error when creating runtime with nil start edge, got nil")
	}

	if runtime != nil {
		runtime.Shutdown()
		t.Error("Expected nil runtime when start edge is nil, got non-nil")
	}

	expectedErrMsg := "runtime creation failed: start edge cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestRuntimeFactory_NilStartNode tests that creating runtime with nil start node fails
func TestRuntimeFactory_NilStartNode(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	initialState := RuntimeTestState{Value: "initial"}

	// Create an edge with nil From node
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: nil, to: node1, role: g.StartEdge}

	runtime, err := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: initialState})
	if err == nil {
		t.Fatal("Expected error when creating runtime with nil start node, got nil")
	}

	if runtime != nil {
		runtime.Shutdown()
		t.Error("Expected nil runtime when start node is nil, got non-nil")
	}

	expectedErrMsg := "runtime creation failed: start node cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestRuntimeFactory_NilTargetNode tests that creating runtime with nil target node fails
func TestRuntimeFactory_NilTargetNode(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)
	initialState := RuntimeTestState{Value: "initial"}

	// Create an edge with nil To node
	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: nil, role: g.StartEdge}

	runtime, err := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: initialState})
	if err == nil {
		t.Fatal("Expected error when creating runtime with nil target node, got nil")
	}

	if runtime != nil {
		runtime.Shutdown()
		t.Error("Expected nil runtime when target node is nil, got non-nil")
	}

	expectedErrMsg := "runtime creation failed: end node cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestRuntimeFactory_NilOptions tests that creating runtime with nil options fails
func TestRuntimeFactory_NilOptions(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, err := RuntimeFactory[RuntimeTestState](startEdge, stateMonitorCh, nil)
	if err == nil {
		t.Fatal("Expected error when creating runtime with nil options, got nil")
	}

	if runtime != nil {
		runtime.Shutdown()
		t.Error("Expected nil runtime when options are nil, got non-nil")
	}

	expectedErrMsg := "runtime creation failed: runtime options cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestRuntime_AddEdge tests adding edges to the runtime
func TestRuntime_AddEdge(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	node2 := newMockRuntimeNode("Node2", g.IntermediateNode, nil, nil)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	edge1 := &mockRuntimeEdge{from: node1, to: node2, role: g.IntermediateEdge}
	edge2 := &mockRuntimeEdge{from: node2, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	// Add edges one at a time
	runtime.AddEdge(edge1)
	runtime.AddEdge(edge2)

	// Verify edges were added (indirectly through validation)
	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate() failed after adding edges: %v", err)
	}
}

// TestRuntime_AddMultipleEdgesAtOnce tests adding multiple edges in one call
func TestRuntime_AddMultipleEdgesAtOnce(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	node2 := newMockRuntimeNode("Node2", g.IntermediateNode, nil, nil)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	edge1 := &mockRuntimeEdge{from: node1, to: node2, role: g.IntermediateEdge}
	edge2 := &mockRuntimeEdge{from: node2, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	// Add multiple edges at once
	runtime.AddEdge(edge1, edge2)

	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate() failed after adding multiple edges: %v", err)
	}
}

// TestRuntime_Validate_ValidGraph tests validation of a valid graph
func TestRuntime_Validate_ValidGraph(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)

	err := runtime.Validate()
	if err != nil {
		t.Errorf("Validate() failed for valid graph: %v", err)
	}
}

// TestRuntime_Validate_NoPathToEnd tests validation failure when no path to end exists
func TestRuntime_Validate_NoPathToEnd(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	// No end edge added
	err := runtime.Validate()
	if err == nil {
		t.Fatal("Expected validation error when no path to end exists, got nil")
	}

	expectedMsg := "graph validation failed: no path from start edge to any end edge"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestRuntime_Invoke_SimpleExecution tests basic graph execution
func TestRuntime_Invoke_SimpleExecution(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		currentState.Value = "processed"
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Counter: 0, Value: "initial"}})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)

	userInput := RuntimeTestState{Value: "input"}
	runtime.Invoke(userInput)

	// Collect state monitor entries
	var entries []g.StateMonitorEntry[RuntimeTestState]
	timeout := time.After(2 * time.Second)

	for {
		select {
		case entry := <-stateMonitorCh:
			entries = append(entries, entry)
			if !entry.Running {
				goto done
			}
		case <-timeout:
			t.Fatal("Test timed out waiting for execution to complete")
		}
	}

done:
	if len(entries) == 0 {
		t.Fatal("No state monitor entries received")
	}

	// Check final state
	finalEntry := entries[len(entries)-1]
	if finalEntry.NewState.Counter != 1 {
		t.Errorf("Expected Counter=1, got %d", finalEntry.NewState.Counter)
	}
	if finalEntry.NewState.Value != "processed" {
		t.Errorf("Expected Value='processed', got '%s'", finalEntry.NewState.Value)
	}
}

// TestRuntime_Invoke_WithError tests execution when a node returns an error
func TestRuntime_Invoke_WithError(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		return currentState, errors.New("node execution failed")
	}, policy)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.Invoke(RuntimeTestState{})

	// Wait for error entry (may get multiple entries before error)
	foundError := false
	timeout := time.After(2 * time.Second)
	for !foundError {
		select {
		case entry := <-stateMonitorCh:
			if entry.Error != nil {
				foundError = true
				if entry.Running {
					t.Error("Expected execution to stop after error, but still running")
				}
			}
		case <-timeout:
			t.Fatal("Test timed out waiting for error entry")
		}
	}
}

// TestRuntime_Invoke_ConcurrentInvocations tests that concurrent invocations are prevented
func TestRuntime_Invoke_ConcurrentInvocations(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		time.Sleep(100 * time.Millisecond) // Simulate long-running task
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)

	// Use the same thread ID for both invocations to test concurrency prevention
	threadID := "test-thread-1"

	// First invocation
	runtime.Invoke(RuntimeTestState{Value: "first"}, g.InvokeConfigThreadID(threadID))

	// Try concurrent invocation on same thread (should fail)
	time.Sleep(10 * time.Millisecond) // Give first invocation time to start
	runtime.Invoke(RuntimeTestState{Value: "second"}, g.InvokeConfigThreadID(threadID))

	// Collect entries
	errorCount := 0
	timeout := time.After(2 * time.Second)

	for {
		select {
		case entry := <-stateMonitorCh:
			if entry.Error != nil && entry.Node == "Runtime" {
				errorCount++
			}
			if !entry.Running && entry.Error == nil {
				goto done
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}

done:
	if errorCount != 1 {
		t.Errorf("Expected 1 concurrent invocation error, got %d", errorCount)
	}
}

// TestRuntime_CurrentState tests retrieving current state
func TestRuntime_CurrentState(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Value: "initial", Counter: 42}})
	defer runtime.Shutdown()

	// Cast to internal implementation to access InitialState
	initialState := runtime.InitialState()
	if initialState.Value != "initial" {
		t.Errorf("Expected Value='initial', got '%s'", initialState.Value)
	}
	if initialState.Counter != 42 {
		t.Errorf("Expected Counter=42, got %d", initialState.Counter)
	}
}

var _ g.Memory[RuntimeTestState] = (*testMemorySetPersistentState)(nil)

type testMemorySetPersistentState struct {
}

func (m *testMemorySetPersistentState) PersistFn() g.PersistFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string, state RuntimeTestState) error {
		return nil
	}
}

func (m *testMemorySetPersistentState) RestoreFn() g.RestoreFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string) (RuntimeTestState, error) {
		return RuntimeTestState{Value: "restored"}, nil
	}
}

// TestRuntime_SetPersistentState tests setting persistence functions
func TestRuntime_SetPersistentState(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{
		Memory: &testMemorySetPersistentState{},
	})
	defer runtime.Shutdown()
	threadID := uuid.NewString()

	// Test that restore works
	err := runtime.Restore(threadID)
	if err != nil {
		t.Errorf("Restore() failed: %v", err)
	}

	restoredState := runtime.CurrentState(threadID)
	if restoredState.Value != "restored" {
		t.Errorf("Expected restored Value='restored', got '%s'", restoredState.Value)
	}
}

// TestRuntime_Restore_WithoutPersistentState tests restore fails without setup
func TestRuntime_Restore_WithoutPersistentState(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	err := runtime.Restore(uuid.NewString())
	if err != nil {
		t.Errorf("Expected no error when restoring without persistence setup, got: %v", err)
	}
}

var _ g.Memory[RuntimeTestState] = (*testMemoryPersistenceStateIsPersisted)(nil)

type testMemoryPersistenceStateIsPersisted struct {
	persistedStates []RuntimeTestState
	mu              sync.Mutex
}

func (m *testMemoryPersistenceStateIsPersisted) PersistFn() g.PersistFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string, state RuntimeTestState) error {
		m.mu.Lock()
		m.persistedStates = append(m.persistedStates, state)
		m.mu.Unlock()
		return nil
	}
}

func (m *testMemoryPersistenceStateIsPersisted) RestoreFn() g.RestoreFn[RuntimeTestState] {
	return func(ctx context.Context, threadID string) (RuntimeTestState, error) {
		return RuntimeTestState{}, nil
	}
}

// TestRuntime_Persistence_StateIsPersisted tests that state changes are persisted
func TestRuntime_Persistence_StateIsPersisted(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 100
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	memory := &testMemoryPersistenceStateIsPersisted{
		persistedStates: make([]RuntimeTestState, 0),
		mu:              sync.Mutex{},
	}

	runtime, _ := RuntimeFactory(
		startEdge,
		stateMonitorCh,
		&g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Counter: 0}, Memory: memory},
	)
	defer runtime.Shutdown()
	runtime.AddEdge(endEdge)

	runtime.Invoke(RuntimeTestState{})

	// Wait for completion
	timeout := time.After(2 * time.Second)
	for {
		select {
		case entry := <-stateMonitorCh:
			if !entry.Running {
				goto done
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}

done:
	// Give persistence worker time to process
	time.Sleep(100 * time.Millisecond)

	memory.mu.Lock()
	count := len(memory.persistedStates)
	memory.mu.Unlock()

	if count == 0 {
		t.Error("Expected at least one persisted state, got 0")
	}
}

// TestRuntime_PartialStateUpdates tests that partial updates are sent to monitor channel
func TestRuntime_PartialStateUpdates(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := &mockRuntimeNode{
		name: "Node1",
		role: g.IntermediateNode,
		fn: func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
			// Emit partial updates
			notify(RuntimeTestState{Value: "partial1", Counter: 1})
			notify(RuntimeTestState{Value: "partial2", Counter: 2})

			// Return final state
			return RuntimeTestState{Value: "final", Counter: 3}, nil
		},
		policy:  policy,
		mailbox: make(chan RuntimeTestState, 10),
	}
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)
	runtime.Invoke(RuntimeTestState{})

	// Collect entries
	var entries []g.StateMonitorEntry[RuntimeTestState]
	timeout := time.After(2 * time.Second)

	for {
		select {
		case entry := <-stateMonitorCh:
			entries = append(entries, entry)
			if !entry.Running {
				goto done
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}

done:
	// Count partial updates
	partialCount := 0
	for _, entry := range entries {
		if entry.Partial {
			partialCount++
		}
	}

	if partialCount != 2 {
		t.Errorf("Expected 2 partial updates, got %d", partialCount)
	}
}

// TestRuntime_MultipleNodes tests execution through multiple nodes
func TestRuntime_MultipleNodes(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}, policy)
	node2 := newMockRuntimeNode("Node2", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	edge1 := &mockRuntimeEdge{from: node1, to: node2, role: g.IntermediateEdge}
	edge2 := &mockRuntimeEdge{from: node2, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Counter: 0}})
	defer runtime.Shutdown()

	runtime.AddEdge(edge1, edge2)
	runtime.Invoke(RuntimeTestState{})

	// Wait for completion
	timeout := time.After(2 * time.Second)
	for {
		select {
		case entry := <-stateMonitorCh:
			if !entry.Running {
				if entry.NewState.Counter != 2 {
					t.Errorf("Expected Counter=2 after two nodes, got %d", entry.NewState.Counter)
				}
				return
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}

// TestRuntime_ConditionalRouting tests routing based on state
func TestRuntime_ConditionalRouting(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	anyPolicy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	conditionalPolicy, _ := RouterPolicyImplFactory(func(userInput, currentState RuntimeTestState, edges []g.Edge[RuntimeTestState]) g.Edge[RuntimeTestState] {
		// Select edge based on state value
		for _, edge := range edges {
			if label, ok := edge.LabelByKey("route"); ok && label == currentState.Value {
				return edge
			}
		}
		return nil
	})

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, anyPolicy)
	routerNode := newMockRuntimeNode("RouterNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Value = "go_left"
		return currentState, nil
	}, conditionalPolicy)
	leftNode := newMockRuntimeNode("LeftNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 100
		return currentState, nil
	}, anyPolicy)
	rightNode := newMockRuntimeNode("RightNode", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 200
		return currentState, nil
	}, anyPolicy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: routerNode, role: g.StartEdge}
	leftEdge := &mockRuntimeEdge{from: routerNode, to: leftNode, role: g.IntermediateEdge, labels: map[string]string{"route": "go_left"}}
	rightEdge := &mockRuntimeEdge{from: routerNode, to: rightNode, role: g.IntermediateEdge, labels: map[string]string{"route": "go_right"}}
	endEdgeLeft := &mockRuntimeEdge{from: leftNode, to: endNode, role: g.EndEdge}
	endEdgeRight := &mockRuntimeEdge{from: rightNode, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Counter: 0}})
	defer runtime.Shutdown()

	runtime.AddEdge(leftEdge, rightEdge, endEdgeLeft, endEdgeRight)
	runtime.Invoke(RuntimeTestState{})

	// Wait for completion
	timeout := time.After(2 * time.Second)
	for {
		select {
		case entry := <-stateMonitorCh:
			if !entry.Running {
				// Should have taken left path, so Counter should be 100
				if entry.NewState.Counter != 100 {
					t.Errorf("Expected Counter=100 (left path), got %d", entry.NewState.Counter)
				}
				return
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}

// TestRuntime_Shutdown tests graceful shutdown
func TestRuntime_Shutdown(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, nil)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, nil, nil)
	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})

	// Shutdown should not panic
	runtime.Shutdown()

	// Calling shutdown again should not panic
	runtime.Shutdown()
}

// TestRuntime_NoOutboundEdges tests error when node has no outbound edges
func TestRuntime_NoOutboundEdges(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		return currentState, nil
	}, policy)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	// Don't add any outbound edges from node1
	runtime.Invoke(RuntimeTestState{})

	// Wait for error (may get multiple entries)
	foundError := false
	timeout := time.After(2 * time.Second)
	for !foundError {
		select {
		case entry := <-stateMonitorCh:
			if entry.Error != nil {
				foundError = true
				expectedMsg := fmt.Sprintf("routing error for node %s: no outbound edges from node", node1.Name())
				if entry.Error.Error() != expectedMsg {
					t.Errorf("Expected error '%s', got '%s'", expectedMsg, entry.Error.Error())
				}
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}

// TestRuntime_NilRoutingPolicy tests error when node has nil routing policy
func TestRuntime_NilRoutingPolicy(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	anyPolicy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, anyPolicy)
	node1 := &mockRuntimeNode{
		name: "Node1",
		role: g.IntermediateNode,
		fn: func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
			return currentState, nil
		},
		policy:  nil, // Nil policy
		mailbox: make(chan RuntimeTestState, 10),
	}
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)
	runtime.Invoke(RuntimeTestState{})

	// Wait for error (may get multiple entries)
	foundError := false
	timeout := time.After(2 * time.Second)
	for !foundError {
		select {
		case entry := <-stateMonitorCh:
			if entry.Error != nil {
				foundError = true
				expectedMsg := fmt.Sprintf("routing error for node %s: no routing policy defined for node", node1.Name())
				if entry.Error.Error() != expectedMsg {
					t.Errorf("Expected error '%s', got '%s'", expectedMsg, entry.Error.Error())
				}
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}

// TestRuntime_EmptyStateMonitorChannel tests runtime works without state monitor channel
func TestRuntime_EmptyStateMonitorChannel(t *testing.T) {
	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter = 42
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	endEdge := &mockRuntimeEdge{from: node1, to: endNode, role: g.EndEdge}

	// Create runtime without state monitor channel (nil)
	runtime, _ := RuntimeFactory[RuntimeTestState](startEdge, nil, &g.RuntimeOptions[RuntimeTestState]{})
	defer runtime.Shutdown()

	runtime.AddEdge(endEdge)
	threadID := runtime.Invoke(RuntimeTestState{})

	// Wait a bit for execution to complete
	time.Sleep(200 * time.Millisecond)

	// Check final state was updated
	finalState := runtime.CurrentState(threadID)
	if finalState.Counter != 42 {
		t.Errorf("Expected Counter=42, got %d", finalState.Counter)
	}
}

// TestRuntime_NonPersistent_CompleteExecution tests runtime execution without any persistence setup
func TestRuntime_NonPersistent_CompleteExecution(t *testing.T) {
	stateMonitorCh := make(chan g.StateMonitorEntry[RuntimeTestState], 10)

	policy, _ := RouterPolicyImplFactory(AnyRoute[RuntimeTestState])

	startNode := newMockRuntimeNode("StartNode", g.StartNode, nil, policy)
	node1 := newMockRuntimeNode("Node1", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		currentState.Value = "modified"
		return currentState, nil
	}, policy)
	node2 := newMockRuntimeNode("Node2", g.IntermediateNode, func(userInput, currentState RuntimeTestState, notify g.NotifyPartialFn[RuntimeTestState]) (RuntimeTestState, error) {
		currentState.Counter++
		currentState.Value = "final"
		return currentState, nil
	}, policy)
	endNode := newMockRuntimeNode("EndNode", g.EndNode, nil, nil)

	startEdge := &mockRuntimeEdge{from: startNode, to: node1, role: g.StartEdge}
	edge1 := &mockRuntimeEdge{from: node1, to: node2, role: g.IntermediateEdge}
	edge2 := &mockRuntimeEdge{from: node2, to: endNode, role: g.EndEdge}

	runtime, _ := RuntimeFactory(startEdge, stateMonitorCh, &g.RuntimeOptions[RuntimeTestState]{InitialState: RuntimeTestState{Counter: 0, Value: "initial"}})
	defer runtime.Shutdown()

	runtime.AddEdge(edge1, edge2)

	// Execute without any persistence setup
	threadID := runtime.Invoke(RuntimeTestState{Value: "user_input"})

	// Collect state monitor entries
	var entries []g.StateMonitorEntry[RuntimeTestState]
	timeout := time.After(2 * time.Second)

	for {
		select {
		case entry := <-stateMonitorCh:
			entries = append(entries, entry)
			if !entry.Running {
				goto done
			}
		case <-timeout:
			t.Fatal("Test timed out waiting for execution to complete")
		}
	}

done:
	if len(entries) == 0 {
		t.Fatal("No state monitor entries received")
	}

	// Verify execution completed successfully
	finalEntry := entries[len(entries)-1]
	if finalEntry.Error != nil {
		t.Errorf("Expected no error, got: %v", finalEntry.Error)
	}

	// Check final state
	if finalEntry.NewState.Counter != 2 {
		t.Errorf("Expected Counter=2, got %d", finalEntry.NewState.Counter)
	}
	if finalEntry.NewState.Value != "final" {
		t.Errorf("Expected Value='final', got '%s'", finalEntry.NewState.Value)
	}

	// Verify state is accessible via CurrentState
	currentState := runtime.CurrentState(threadID)
	if currentState.Counter != 2 {
		t.Errorf("CurrentState Counter expected 2, got %d", currentState.Counter)
	}
	if currentState.Value != "final" {
		t.Errorf("CurrentState Value expected 'final', got '%s'", currentState.Value)
	}

	// Verify nodes were called
	if node1.GetCallCount() != 1 {
		t.Errorf("Expected node1 to be called once, got %d", node1.GetCallCount())
	}
	if node2.GetCallCount() != 1 {
		t.Errorf("Expected node2 to be called once, got %d", node2.GetCallCount())
	}
}
