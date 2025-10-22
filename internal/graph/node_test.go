package graph_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/morphy76/ggraph/internal/graph"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// NodeTestState is a state type for node testing with counter
type NodeTestState struct {
	Value   string
	Counter int
}

// mockStateObserver is a minimal StateObserver implementation for testing
type mockStateObserver struct {
	mu              sync.Mutex
	currentState    NodeTestState
	notifications   []stateNotification
	notificationsCh chan stateNotification
}

type stateNotification struct {
	nodeName    string
	userInput   NodeTestState
	stateChange NodeTestState
	err         error
	partial     bool
}

func newMockStateObserver(initialState NodeTestState) *mockStateObserver {
	return &mockStateObserver{
		currentState:    initialState,
		notifications:   make([]stateNotification, 0),
		notificationsCh: make(chan stateNotification, 10),
	}
}

func (m *mockStateObserver) NotifyStateChange(node g.Node[NodeTestState], userInput, stateChange NodeTestState, reducer g.ReducerFn[NodeTestState], err error, partial bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if reducer != nil {
		m.currentState = reducer(m.currentState, stateChange)
	} else {
		m.currentState = stateChange
	}

	notification := stateNotification{
		nodeName:    node.Name(),
		userInput:   userInput,
		stateChange: stateChange,
		err:         err,
		partial:     partial,
	}
	m.notifications = append(m.notifications, notification)
	m.notificationsCh <- notification
}

func (m *mockStateObserver) CurrentState() NodeTestState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentState
}

func (m *mockStateObserver) getNotifications() []stateNotification {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]stateNotification{}, m.notifications...)
}

func TestNodeImplFactory_BasicCreation(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter++
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"test-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	if node == nil {
		t.Fatal("NodeImplFactory returned nil")
	}

	if node.Name() != "test-node" {
		t.Errorf("Expected Name() to return 'test-node', got '%s'", node.Name())
	}

	if node.Role() != g.IntermediateNode {
		t.Errorf("Expected Role() to return IntermediateNode, got %v", node.Role())
	}

	if node.RoutePolicy() == nil {
		t.Error("Expected RoutePolicy() to return non-nil value")
	}
}

func TestNodeImplFactory_WithNilNodeFn(t *testing.T) {
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"test-node",
		nil,
		routePolicy,
		reducer,
	)

	if node == nil {
		t.Fatal("NodeImplFactory returned nil")
	}

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 5}

	node.Accept(userInput, observer)

	select {
	case notification := <-observer.notificationsCh:
		if notification.err != nil {
			t.Errorf("Unexpected error: %v", notification.err)
		}
		if notification.stateChange.Value != "initial" {
			t.Errorf("Expected default function to preserve state, got Value=%s", notification.stateChange.Value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for node execution")
	}
}

func TestNodeImplFactory_WithNilRoutePolicy(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"test-node",
		nodeFn,
		nil,
		reducer,
	)

	if node == nil {
		t.Fatal("NodeImplFactory returned nil")
	}

	if node.RoutePolicy() == nil {
		t.Error("Expected RoutePolicy() to return non-nil default value")
	}
}

func TestNodeImplFactory_AllRoles(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	testCases := []struct {
		name string
		role g.NodeRole
	}{
		{"StartNode", g.StartNode},
		{"IntermediateNode", g.IntermediateNode},
		{"EndNode", g.EndNode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := graph.NodeImplFactory[NodeTestState](
				tc.role,
				"test-node",
				nodeFn,
				routePolicy,
				reducer,
			)

			if node.Role() != tc.role {
				t.Errorf("Expected Role() to return %v, got %v", tc.role, node.Role())
			}
		})
	}
}

func TestNodeImplFactory_NodeExecution(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		currentState.Counter = userInput.Counter + 10
		currentState.Value = "processed"
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"process-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 5}

	node.Accept(userInput, observer)

	select {
	case notification := <-observer.notificationsCh:
		if notification.err != nil {
			t.Errorf("Unexpected error: %v", notification.err)
		}
		if notification.partial {
			t.Error("Expected final notification, got partial")
		}
		if notification.nodeName != "process-node" {
			t.Errorf("Expected nodeName 'process-node', got '%s'", notification.nodeName)
		}
		if notification.stateChange.Counter != 15 {
			t.Errorf("Expected Counter=15, got %d", notification.stateChange.Counter)
		}
		if notification.stateChange.Value != "processed" {
			t.Errorf("Expected Value='processed', got '%s'", notification.stateChange.Value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for node execution")
	}
}

func TestNodeImplFactory_NodeExecutionWithError(t *testing.T) {
	expectedErr := errors.New("processing failed")
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, expectedErr
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"error-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 5}

	node.Accept(userInput, observer)

	select {
	case notification := <-observer.notificationsCh:
		if notification.err == nil {
			t.Error("Expected error, got nil")
		}
		if notification.partial {
			t.Error("Expected final notification, got partial")
		}
		if notification.err != nil && notification.nodeName != "error-node" {
			t.Errorf("Expected nodeName 'error-node', got '%s'", notification.nodeName)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for node execution")
	}
}

func TestNodeImplFactory_PartialStateUpdates(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		for i := 1; i <= 3; i++ {
			currentState.Counter = i
			notify(currentState)
			time.Sleep(10 * time.Millisecond)
		}
		currentState.Counter = 10
		currentState.Value = "complete"
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"partial-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
	userInput := NodeTestState{Value: "input", Counter: 0}

	node.Accept(userInput, observer)

	notifications := make([]stateNotification, 0)
	timeout := time.After(2 * time.Second)

	for i := 0; i < 4; i++ {
		select {
		case notification := <-observer.notificationsCh:
			notifications = append(notifications, notification)
		case <-timeout:
			t.Fatalf("Timeout waiting for notification %d", i+1)
		}
	}

	partialCount := 0
	finalCount := 0
	for _, n := range notifications {
		if n.partial {
			partialCount++
		} else {
			finalCount++
		}
	}

	if partialCount != 3 {
		t.Errorf("Expected 3 partial notifications, got %d", partialCount)
	}
	if finalCount != 1 {
		t.Errorf("Expected 1 final notification, got %d", finalCount)
	}

	finalNotification := notifications[len(notifications)-1]
	if finalNotification.partial {
		t.Error("Last notification should not be partial")
	}
	if finalNotification.stateChange.Counter != 10 {
		t.Errorf("Expected final Counter=10, got %d", finalNotification.stateChange.Counter)
	}
	if finalNotification.stateChange.Value != "complete" {
		t.Errorf("Expected final Value='complete', got '%s'", finalNotification.stateChange.Value)
	}
}

func TestNodeImplFactory_DifferentStateTypes(t *testing.T) {
	type AnotherState struct {
		Message string
		Count   int
	}

	nodeFn := func(userInput, currentState AnotherState, notify g.NotifyPartialFn[AnotherState]) (AnotherState, error) {
		currentState.Count++
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[AnotherState](graph.AnyRoute[AnotherState])
	reducer := graph.Replacer[AnotherState]

	node := graph.NodeImplFactory[AnotherState](
		g.IntermediateNode,
		"typed-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	if node == nil {
		t.Error("NodeImplFactory failed to create node with AnotherState")
	}

	if node.Name() != "typed-node" {
		t.Errorf("Expected Name() to return 'typed-node', got '%s'", node.Name())
	}
}

func TestNodeImplFactory_WithReducer(t *testing.T) {
	customReducer := func(currentState, change NodeTestState) NodeTestState {
		currentState.Counter += change.Counter
		if change.Value != "" {
			currentState.Value = change.Value
		}
		return currentState
	}

	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return NodeTestState{Counter: 5, Value: "added"}, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"reducer-node",
		nodeFn,
		routePolicy,
		customReducer,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 10})
	userInput := NodeTestState{Value: "input", Counter: 0}

	node.Accept(userInput, observer)

	select {
	case notification := <-observer.notificationsCh:
		if notification.err != nil {
			t.Errorf("Unexpected error: %v", notification.err)
		}
		finalState := observer.CurrentState()
		if finalState.Counter != 15 {
			t.Errorf("Expected Counter=15 (10+5), got %d", finalState.Counter)
		}
		if finalState.Value != "added" {
			t.Errorf("Expected Value='added', got '%s'", finalState.Value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for node execution")
	}
}

func TestNodeImplFactory_MultipleExecutions(t *testing.T) {
	executionCount := 0
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		executionCount++
		currentState.Counter = executionCount
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"multi-exec-node",
		nodeFn,
		routePolicy,
		reducer,
	)

	for i := 1; i <= 3; i++ {
		observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 0})
		userInput := NodeTestState{Value: "input", Counter: i}

		node.Accept(userInput, observer)

		select {
		case notification := <-observer.notificationsCh:
			if notification.err != nil {
				t.Errorf("Execution %d: Unexpected error: %v", i, notification.err)
			}
			if notification.stateChange.Counter != i {
				t.Errorf("Execution %d: Expected Counter=%d, got %d", i, i, notification.stateChange.Counter)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("Execution %d: Timeout waiting for node execution", i)
		}
	}

	if executionCount != 3 {
		t.Errorf("Expected 3 executions, got %d", executionCount)
	}
}

func TestNodeImplFactory_NodeNamePreservation(t *testing.T) {
	names := []string{"node-1", "processing-unit", "validator", "transformer"}

	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return currentState, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])
	reducer := graph.Replacer[NodeTestState]

	for _, name := range names {
		node := graph.NodeImplFactory[NodeTestState](
			g.IntermediateNode,
			name,
			nodeFn,
			routePolicy,
			reducer,
		)

		if node.Name() != name {
			t.Errorf("Expected Name() to return '%s', got '%s'", name, node.Name())
		}
	}
}

func TestNodeImplFactory_NilReducer(t *testing.T) {
	nodeFn := func(userInput, currentState NodeTestState, notify g.NotifyPartialFn[NodeTestState]) (NodeTestState, error) {
		return NodeTestState{Counter: 99, Value: "new"}, nil
	}
	routePolicy, _ := graph.RouterPolicyImplFactory[NodeTestState](graph.AnyRoute[NodeTestState])

	node := graph.NodeImplFactory[NodeTestState](
		g.IntermediateNode,
		"nil-reducer-node",
		nodeFn,
		routePolicy,
		nil,
	)

	observer := newMockStateObserver(NodeTestState{Value: "initial", Counter: 10})
	userInput := NodeTestState{Value: "input", Counter: 0}

	node.Accept(userInput, observer)

	select {
	case notification := <-observer.notificationsCh:
		if notification.err != nil {
			t.Errorf("Unexpected error: %v", notification.err)
		}
		finalState := observer.CurrentState()
		if finalState.Counter != 99 {
			t.Errorf("Expected Counter=99, got %d", finalState.Counter)
		}
		if finalState.Value != "new" {
			t.Errorf("Expected Value='new', got '%s'", finalState.Value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for node execution")
	}
}
