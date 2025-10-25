package graph

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// StateTestState is a state type for state monitoring testing
type StateTestState struct {
	Value   string
	Counter int
	Data    map[string]interface{}
}

// Test monitorRunning function

func TestMonitorRunning_BasicState(t *testing.T) {
	state := StateTestState{
		Value:   "test",
		Counter: 42,
	}

	threadID := uuid.NewString()

	entry := monitorRunning("test-node", threadID, state)

	if entry.Node != "test-node" {
		t.Errorf("Expected Node='test-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if entry.ThreadID == "" {
		t.Error("Expected ThreadID to be non-empty")
	}

	if !entry.Running {
		t.Error("Expected Running=true")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.Error != nil {
		t.Errorf("Expected Error=nil, got %v", entry.Error)
	}

	if entry.NewState.Value != "test" {
		t.Errorf("Expected NewState.Value='test', got '%s'", entry.NewState.Value)
	}

	if entry.NewState.Counter != 42 {
		t.Errorf("Expected NewState.Counter=42, got %d", entry.NewState.Counter)
	}
}

func TestMonitorRunning_EmptyState(t *testing.T) {
	state := StateTestState{}

	threadID := uuid.NewString()

	entry := monitorRunning("empty-node", threadID, state)

	if entry.Node != "empty-node" {
		t.Errorf("Expected Node='empty-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if !entry.Running {
		t.Error("Expected Running=true")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.NewState.Value != "" {
		t.Errorf("Expected empty NewState.Value, got '%s'", entry.NewState.Value)
	}
}

func TestMonitorRunning_ComplexState(t *testing.T) {
	state := StateTestState{
		Value:   "complex",
		Counter: 100,
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	threadID := uuid.NewString()

	entry := monitorRunning("complex-node", threadID, state)

	if entry.Node != "complex-node" {
		t.Errorf("Expected Node='complex-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if !entry.Running {
		t.Error("Expected Running=true")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.NewState.Data == nil {
		t.Fatal("Expected NewState.Data to be non-nil")
	}

	if len(entry.NewState.Data) != 2 {
		t.Errorf("Expected NewState.Data length=2, got %d", len(entry.NewState.Data))
	}
}

// Test monitorNonFatalError function

func TestMonitorNonFatalError_WithError(t *testing.T) {
	testErr := errors.New("non-fatal error occurred")

	threadID := uuid.NewString()

	entry := monitorNonFatalError[StateTestState]("error-node", threadID, testErr)

	if entry.Node != "error-node" {
		t.Errorf("Expected Node='error-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if !entry.Running {
		t.Error("Expected Running=true for non-fatal error")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.Error == nil {
		t.Fatal("Expected Error to be non-nil")
	}

	if entry.Error.Error() != "non-fatal error occurred" {
		t.Errorf("Expected error message 'non-fatal error occurred', got '%s'", entry.Error.Error())
	}

	// NewState should be zero value
	if entry.NewState.Value != "" || entry.NewState.Counter != 0 {
		t.Error("Expected NewState to be zero value for error entry")
	}
}

func TestMonitorNonFatalError_DifferentErrors(t *testing.T) {
	errors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("timeout"),
	}

	threadID := uuid.NewString()

	for i, err := range errors {
		entry := monitorNonFatalError[StateTestState]("node", threadID, err)

		if entry.ThreadID != threadID {
			t.Errorf("Test %d: Expected ThreadID='%s', got '%s'", i, threadID, entry.ThreadID)
		}

		if !entry.Running {
			t.Errorf("Test %d: Expected Running=true", i)
		}

		if entry.Error != err {
			t.Errorf("Test %d: Expected Error=%v, got %v", i, err, entry.Error)
		}
	}
}

// Test monitorError function

func TestMonitorError_WithError(t *testing.T) {
	testErr := errors.New("fatal error occurred")

	threadID := uuid.NewString()

	entry := monitorError[StateTestState]("fatal-node", threadID, testErr)

	if entry.Node != "fatal-node" {
		t.Errorf("Expected Node='fatal-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if entry.Running {
		t.Error("Expected Running=false for fatal error")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.Error == nil {
		t.Fatal("Expected Error to be non-nil")
	}

	if entry.Error.Error() != "fatal error occurred" {
		t.Errorf("Expected error message 'fatal error occurred', got '%s'", entry.Error.Error())
	}
}

func TestMonitorError_StopsExecution(t *testing.T) {
	testErr := errors.New("stop execution")

	threadID := uuid.NewString()

	entry := monitorError[StateTestState]("stop-node", threadID, testErr)

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	// The key difference from non-fatal is Running=false
	if entry.Running {
		t.Error("Expected Running=false, execution should stop")
	}
}

func TestMonitorError_VsNonFatalError(t *testing.T) {
	testErr := errors.New("test error")

	threadID := uuid.NewString()

	nonFatalEntry := monitorNonFatalError[StateTestState]("node", threadID, testErr)
	fatalEntry := monitorError[StateTestState]("node", threadID, testErr)

	// Both should have the same thread ID
	if nonFatalEntry.ThreadID != threadID {
		t.Errorf("Expected non-fatal ThreadID='%s', got '%s'", threadID, nonFatalEntry.ThreadID)
	}
	if fatalEntry.ThreadID != threadID {
		t.Errorf("Expected fatal ThreadID='%s', got '%s'", threadID, fatalEntry.ThreadID)
	}

	// Both should have the same error
	if nonFatalEntry.Error != fatalEntry.Error {
		t.Error("Both entries should have the same error")
	}

	// But Running should be different
	if !nonFatalEntry.Running {
		t.Error("Expected non-fatal error to have Running=true")
	}

	if fatalEntry.Running {
		t.Error("Expected fatal error to have Running=false")
	}
}

// Test monitorPartial function

func TestMonitorPartial_BasicState(t *testing.T) {
	state := StateTestState{
		Value:   "partial",
		Counter: 5,
	}

	threadID := uuid.NewString()

	entry := monitorPartial("partial-node", threadID, state)

	if entry.Node != "partial-node" {
		t.Errorf("Expected Node='partial-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if !entry.Running {
		t.Error("Expected Running=true")
	}

	if !entry.Partial {
		t.Error("Expected Partial=true")
	}

	if entry.Error != nil {
		t.Errorf("Expected Error=nil, got %v", entry.Error)
	}

	if entry.NewState.Value != "partial" {
		t.Errorf("Expected NewState.Value='partial', got '%s'", entry.NewState.Value)
	}
}

func TestMonitorPartial_MultipleUpdates(t *testing.T) {
	states := []StateTestState{
		{Value: "update1", Counter: 1},
		{Value: "update2", Counter: 2},
		{Value: "update3", Counter: 3},
	}

	threadID := uuid.NewString()

	for i, state := range states {
		entry := monitorPartial("node", threadID, state)

		if entry.ThreadID != threadID {
			t.Errorf("Update %d: Expected ThreadID='%s', got '%s'", i, threadID, entry.ThreadID)
		}

		if !entry.Partial {
			t.Errorf("Update %d: Expected Partial=true", i)
		}

		if !entry.Running {
			t.Errorf("Update %d: Expected Running=true", i)
		}

		if entry.NewState.Counter != i+1 {
			t.Errorf("Update %d: Expected Counter=%d, got %d", i, i+1, entry.NewState.Counter)
		}
	}
}

func TestMonitorPartial_VsRunning(t *testing.T) {
	state := StateTestState{Value: "test", Counter: 10}

	threadID := uuid.NewString()

	partialEntry := monitorPartial("node", threadID, state)
	runningEntry := monitorRunning("node", threadID, state)

	// Both should have the same thread ID
	if partialEntry.ThreadID != threadID || runningEntry.ThreadID != threadID {
		t.Errorf("Both entries should have ThreadID='%s'", threadID)
	}

	// Both should be running
	if !partialEntry.Running || !runningEntry.Running {
		t.Error("Both entries should have Running=true")
	}

	// But Partial should be different
	if !partialEntry.Partial {
		t.Error("Expected partial entry to have Partial=true")
	}

	if runningEntry.Partial {
		t.Error("Expected running entry to have Partial=false")
	}

	// Both should have the same state values
	if partialEntry.NewState.Value != runningEntry.NewState.Value {
		t.Error("Both entries should have the same NewState.Value")
	}
	if partialEntry.NewState.Counter != runningEntry.NewState.Counter {
		t.Error("Both entries should have the same NewState.Counter")
	}
}

// Test monitorCompleted function

func TestMonitorCompleted_BasicState(t *testing.T) {
	state := StateTestState{
		Value:   "completed",
		Counter: 100,
	}

	threadID := uuid.NewString()

	entry := monitorCompleted("completed-node", threadID, state)

	if entry.Node != "completed-node" {
		t.Errorf("Expected Node='completed-node', got '%s'", entry.Node)
	}

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if entry.Running {
		t.Error("Expected Running=false")
	}

	if entry.Partial {
		t.Error("Expected Partial=false")
	}

	if entry.Error != nil {
		t.Errorf("Expected Error=nil, got %v", entry.Error)
	}

	if entry.NewState.Value != "completed" {
		t.Errorf("Expected NewState.Value='completed', got '%s'", entry.NewState.Value)
	}

	if entry.NewState.Counter != 100 {
		t.Errorf("Expected NewState.Counter=100, got %d", entry.NewState.Counter)
	}
}

func TestMonitorCompleted_FinalState(t *testing.T) {
	state := StateTestState{
		Value:   "final",
		Counter: 999,
		Data: map[string]interface{}{
			"result": "success",
		},
	}

	threadID := uuid.NewString()

	entry := monitorCompleted("final-node", threadID, state)

	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}

	if entry.Running {
		t.Error("Expected Running=false, execution should be complete")
	}

	if entry.Partial {
		t.Error("Expected Partial=false, this is the final state")
	}

	if entry.Error != nil {
		t.Error("Expected Error=nil for successful completion")
	}
}

func TestMonitorCompleted_VsRunning(t *testing.T) {
	state := StateTestState{Value: "test", Counter: 50}

	threadID := uuid.NewString()

	completedEntry := monitorCompleted("node", threadID, state)
	runningEntry := monitorRunning("node", threadID, state)

	// Both should have the same thread ID
	if completedEntry.ThreadID != threadID || runningEntry.ThreadID != threadID {
		t.Errorf("Both entries should have ThreadID='%s'", threadID)
	}

	// Running should be different
	if completedEntry.Running {
		t.Error("Expected completed entry to have Running=false")
	}

	if !runningEntry.Running {
		t.Error("Expected running entry to have Running=true")
	}

	// Both should have same state values
	if completedEntry.NewState.Value != runningEntry.NewState.Value {
		t.Error("Both entries should have the same NewState.Value")
	}
	if completedEntry.NewState.Counter != runningEntry.NewState.Counter {
		t.Error("Both entries should have the same NewState.Counter")
	}

	// Both should be non-partial
	if completedEntry.Partial || runningEntry.Partial {
		t.Error("Neither should be partial")
	}
}

// Test combinations and edge cases

func TestMonitorFunctions_AllNodeNames(t *testing.T) {
	state := StateTestState{Value: "test"}
	nodeNames := []string{"node1", "processing-unit", "validator", "end"}

	threadID := uuid.NewString()

	for _, nodeName := range nodeNames {
		runningEntry := monitorRunning(nodeName, threadID, state)
		if runningEntry.Node != nodeName {
			t.Errorf("MonitorRunning: expected Node='%s', got '%s'", nodeName, runningEntry.Node)
		}
		if runningEntry.ThreadID != threadID {
			t.Errorf("MonitorRunning: expected ThreadID='%s', got '%s'", threadID, runningEntry.ThreadID)
		}

		partialEntry := monitorPartial(nodeName, threadID, state)
		if partialEntry.Node != nodeName {
			t.Errorf("MonitorPartial: expected Node='%s', got '%s'", nodeName, partialEntry.Node)
		}
		if partialEntry.ThreadID != threadID {
			t.Errorf("MonitorPartial: expected ThreadID='%s', got '%s'", threadID, partialEntry.ThreadID)
		}

		completedEntry := monitorCompleted(nodeName, threadID, state)
		if completedEntry.Node != nodeName {
			t.Errorf("MonitorCompleted: expected Node='%s', got '%s'", nodeName, completedEntry.Node)
		}
		if completedEntry.ThreadID != threadID {
			t.Errorf("MonitorCompleted: expected ThreadID='%s', got '%s'", threadID, completedEntry.ThreadID)
		}
	}
}

func TestMonitorFunctions_CompareRunningAndPartialFlags(t *testing.T) {
	state := StateTestState{Value: "test", Counter: 10}
	node := "test-node"

	threadID := uuid.NewString()

	testCases := []struct {
		name             string
		createEntry      func() g.StateMonitorEntry[StateTestState]
		expectedRunning  bool
		expectedPartial  bool
		expectedHasError bool
	}{
		{
			name:             "monitorRunning",
			createEntry:      func() g.StateMonitorEntry[StateTestState] { return monitorRunning(node, threadID, state) },
			expectedRunning:  true,
			expectedPartial:  false,
			expectedHasError: false,
		},
		{
			name:             "monitorPartial",
			createEntry:      func() g.StateMonitorEntry[StateTestState] { return monitorPartial(node, threadID, state) },
			expectedRunning:  true,
			expectedPartial:  true,
			expectedHasError: false,
		},
		{
			name:             "monitorCompleted",
			createEntry:      func() g.StateMonitorEntry[StateTestState] { return monitorCompleted(node, threadID, state) },
			expectedRunning:  false,
			expectedPartial:  false,
			expectedHasError: false,
		},
		{
			name: "monitorNonFatalError",
			createEntry: func() g.StateMonitorEntry[StateTestState] {
				return monitorNonFatalError[StateTestState](node, threadID, errors.New("test"))
			},
			expectedRunning:  true,
			expectedPartial:  false,
			expectedHasError: true,
		},
		{
			name: "monitorError",
			createEntry: func() g.StateMonitorEntry[StateTestState] {
				return monitorError[StateTestState](node, threadID, errors.New("test"))
			},
			expectedRunning:  false,
			expectedPartial:  false,
			expectedHasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := tc.createEntry()

			if entry.Running != tc.expectedRunning {
				t.Errorf("Expected Running=%v, got %v", tc.expectedRunning, entry.Running)
			}

			if entry.Partial != tc.expectedPartial {
				t.Errorf("Expected Partial=%v, got %v", tc.expectedPartial, entry.Partial)
			}

			hasError := entry.Error != nil
			if hasError != tc.expectedHasError {
				t.Errorf("Expected hasError=%v, got %v", tc.expectedHasError, hasError)
			}
		})
	}
}

func TestMonitorFunctions_ExecutionLifecycle(t *testing.T) {
	// Simulate a typical execution lifecycle
	nodeName := "processing-node"

	threadID := uuid.NewString()

	// 1. Node starts executing
	runningEntry := monitorRunning(nodeName, threadID, StateTestState{Counter: 0})
	if runningEntry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, runningEntry.ThreadID)
	}
	if !runningEntry.Running || runningEntry.Partial {
		t.Error("Initial execution should have Running=true, Partial=false")
	}

	// 2. Node sends partial updates
	partial1 := monitorPartial(nodeName, threadID, StateTestState{Counter: 1})
	partial2 := monitorPartial(nodeName, threadID, StateTestState{Counter: 2})
	partial3 := monitorPartial(nodeName, threadID, StateTestState{Counter: 3})

	for i, entry := range []g.StateMonitorEntry[StateTestState]{partial1, partial2, partial3} {
		if entry.ThreadID != threadID {
			t.Errorf("Partial update %d: expected ThreadID='%s', got '%s'", i+1, threadID, entry.ThreadID)
		}
		if !entry.Running || !entry.Partial {
			t.Errorf("Partial update %d should have Running=true, Partial=true", i+1)
		}
	}

	// 3. Node completes successfully
	completedEntry := monitorCompleted(nodeName, threadID, StateTestState{Counter: 10})
	if completedEntry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, completedEntry.ThreadID)
	}
	if completedEntry.Running || completedEntry.Partial {
		t.Error("Completion should have Running=false, Partial=false")
	}
}

func TestMonitorFunctions_ErrorScenarios(t *testing.T) {
	nodeName := "error-node"

	threadID := uuid.NewString()

	// Scenario 1: Non-fatal error during execution (can continue)
	nonFatalEntry := monitorNonFatalError[StateTestState](nodeName, threadID, errors.New("retry possible"))
	if nonFatalEntry.ThreadID != threadID {
		t.Errorf("Expected non-fatal ThreadID='%s', got '%s'", threadID, nonFatalEntry.ThreadID)
	}
	if !nonFatalEntry.Running {
		t.Error("Non-fatal error should allow continued execution (Running=true)")
	}

	// Scenario 2: Fatal error stops execution
	fatalEntry := monitorError[StateTestState](nodeName, threadID, errors.New("critical failure"))
	if fatalEntry.ThreadID != threadID {
		t.Errorf("Expected fatal ThreadID='%s', got '%s'", threadID, fatalEntry.ThreadID)
	}
	if fatalEntry.Running {
		t.Error("Fatal error should stop execution (Running=false)")
	}

	// Both should have errors
	if nonFatalEntry.Error == nil || fatalEntry.Error == nil {
		t.Error("Both error scenarios should have Error set")
	}
}

func TestMonitorFunctions_DifferentStateTypes(t *testing.T) {
	type AnotherState struct {
		Name    string
		Active  bool
		Results []int
	}

	state := AnotherState{
		Name:    "test",
		Active:  true,
		Results: []int{1, 2, 3},
	}

	threadID := uuid.NewString()

	runningEntry := monitorRunning("node", threadID, state)
	if runningEntry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, runningEntry.ThreadID)
	}
	if runningEntry.NewState.Name != "test" {
		t.Errorf("Expected Name='test', got '%s'", runningEntry.NewState.Name)
	}

	completedEntry := monitorCompleted("node", threadID, state)
	if completedEntry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, completedEntry.ThreadID)
	}
	if completedEntry.NewState.Active != true {
		t.Error("Expected Active=true")
	}

	partialEntry := monitorPartial("node", threadID, state)
	if partialEntry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, partialEntry.ThreadID)
	}
	if len(partialEntry.NewState.Results) != 3 {
		t.Errorf("Expected Results length=3, got %d", len(partialEntry.NewState.Results))
	}
}

func TestMonitorFunctions_EmptyNodeName(t *testing.T) {
	state := StateTestState{Value: "test"}

	threadID := uuid.NewString()

	// Test with empty node name
	entry := monitorRunning("", threadID, state)
	if entry.Node != "" {
		t.Errorf("Expected empty Node, got '%s'", entry.Node)
	}
	if entry.ThreadID != threadID {
		t.Errorf("Expected ThreadID='%s', got '%s'", threadID, entry.ThreadID)
	}
}

func TestMonitorFunctions_NilError(t *testing.T) {
	threadID := uuid.NewString()

	// These functions should handle nil errors gracefully
	nonFatalEntry := monitorNonFatalError[StateTestState]("node", threadID, nil)
	if nonFatalEntry.ThreadID != threadID {
		t.Errorf("Expected non-fatal ThreadID='%s', got '%s'", threadID, nonFatalEntry.ThreadID)
	}
	if nonFatalEntry.Error != nil {
		t.Error("Expected Error=nil when nil is passed")
	}

	fatalEntry := monitorError[StateTestState]("node", threadID, nil)
	if fatalEntry.ThreadID != threadID {
		t.Errorf("Expected fatal ThreadID='%s', got '%s'", threadID, fatalEntry.ThreadID)
	}
	if fatalEntry.Error != nil {
		t.Error("Expected Error=nil when nil is passed")
	}
}
