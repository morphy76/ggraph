package graph_test

import (
	"testing"

	"github.com/morphy76/ggraph/internal/graph"
)

// SimpleState is a simple state type for testing
type SimpleState struct {
	Value   string
	Counter int
}

// ComplexState is a more complex state type for testing
type ComplexState struct {
	ID       string
	Data     map[string]interface{}
	Items    []string
	Metadata struct {
		Created string
		Updated string
	}
}

func TestReplacer_SimpleState(t *testing.T) {
	currentState := SimpleState{
		Value:   "current",
		Counter: 10,
	}

	changeState := SimpleState{
		Value:   "new",
		Counter: 20,
	}

	result := graph.Replacer(currentState, changeState)

	// Replacer should completely replace currentState with changeState
	if result.Value != changeState.Value {
		t.Errorf("Expected Value=%s, got %s", changeState.Value, result.Value)
	}
	if result.Counter != changeState.Counter {
		t.Errorf("Expected Counter=%d, got %d", changeState.Counter, result.Counter)
	}

	// Verify it's not the currentState
	if result.Value == currentState.Value {
		t.Error("Result should not have currentState values")
	}
	if result.Counter == currentState.Counter {
		t.Error("Result should not have currentState counter")
	}
}

func TestReplacer_EmptyCurrentState(t *testing.T) {
	currentState := SimpleState{}

	changeState := SimpleState{
		Value:   "new",
		Counter: 42,
	}

	result := graph.Replacer(currentState, changeState)

	if result.Value != "new" {
		t.Errorf("Expected Value='new', got '%s'", result.Value)
	}
	if result.Counter != 42 {
		t.Errorf("Expected Counter=42, got %d", result.Counter)
	}
}

func TestReplacer_EmptyChangeState(t *testing.T) {
	currentState := SimpleState{
		Value:   "current",
		Counter: 10,
	}

	changeState := SimpleState{}

	result := graph.Replacer(currentState, changeState)

	// Even with empty change state, it should replace
	if result.Value != "" {
		t.Errorf("Expected Value='', got '%s'", result.Value)
	}
	if result.Counter != 0 {
		t.Errorf("Expected Counter=0, got %d", result.Counter)
	}
}

func TestReplacer_BothEmptyStates(t *testing.T) {
	currentState := SimpleState{}
	changeState := SimpleState{}

	result := graph.Replacer(currentState, changeState)

	if result.Value != "" {
		t.Errorf("Expected Value='', got '%s'", result.Value)
	}
	if result.Counter != 0 {
		t.Errorf("Expected Counter=0, got %d", result.Counter)
	}
}

func TestReplacer_ComplexState(t *testing.T) {
	currentState := ComplexState{
		ID: "old-id",
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
		Items: []string{"item1", "item2"},
		Metadata: struct {
			Created string
			Updated string
		}{
			Created: "2024-01-01",
			Updated: "2024-01-02",
		},
	}

	changeState := ComplexState{
		ID: "new-id",
		Data: map[string]interface{}{
			"key3": "value3",
		},
		Items: []string{"item3"},
		Metadata: struct {
			Created string
			Updated string
		}{
			Created: "2025-01-01",
			Updated: "2025-01-02",
		},
	}

	result := graph.Replacer(currentState, changeState)

	// Verify complete replacement
	if result.ID != "new-id" {
		t.Errorf("Expected ID='new-id', got '%s'", result.ID)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected Data length=1, got %d", len(result.Data))
	}

	if val, ok := result.Data["key3"]; !ok || val != "value3" {
		t.Error("Expected Data to contain key3='value3'")
	}

	if _, ok := result.Data["key1"]; ok {
		t.Error("Expected Data to not contain old key1")
	}

	if len(result.Items) != 1 || result.Items[0] != "item3" {
		t.Errorf("Expected Items=['item3'], got %v", result.Items)
	}

	if result.Metadata.Created != "2025-01-01" {
		t.Errorf("Expected Metadata.Created='2025-01-01', got '%s'", result.Metadata.Created)
	}
}

func TestReplacer_ZeroValues(t *testing.T) {
	currentState := SimpleState{
		Value:   "something",
		Counter: 100,
	}

	changeState := SimpleState{
		Value:   "",
		Counter: 0,
	}

	result := graph.Replacer(currentState, changeState)

	// Replacer should replace with zero values
	if result.Value != "" {
		t.Errorf("Expected Value='', got '%s'", result.Value)
	}
	if result.Counter != 0 {
		t.Errorf("Expected Counter=0, got %d", result.Counter)
	}
}

func TestReplacer_PartialUpdate(t *testing.T) {
	// Note: Replacer does NOT do partial updates, it completely replaces
	currentState := SimpleState{
		Value:   "current",
		Counter: 10,
	}

	changeState := SimpleState{
		Value: "updated",
		// Counter is zero value
	}

	result := graph.Replacer(currentState, changeState)

	// Replacer replaces everything, including zero values
	if result.Value != "updated" {
		t.Errorf("Expected Value='updated', got '%s'", result.Value)
	}
	if result.Counter != 0 {
		t.Errorf("Expected Counter=0 (zero value replacement), got %d", result.Counter)
	}
}

func TestReplacer_DifferentStateTypes(t *testing.T) {
	type StateA struct {
		Name string
		Age  int
	}

	currentA := StateA{Name: "Alice", Age: 30}
	changeA := StateA{Name: "Bob", Age: 25}

	resultA := graph.Replacer(currentA, changeA)

	if resultA.Name != "Bob" {
		t.Errorf("Expected Name='Bob', got '%s'", resultA.Name)
	}
	if resultA.Age != 25 {
		t.Errorf("Expected Age=25, got %d", resultA.Age)
	}

	type StateB struct {
		Title       string
		Count       int
		IsActive    bool
		Temperature float64
	}

	currentB := StateB{Title: "Test", Count: 5, IsActive: true, Temperature: 20.5}
	changeB := StateB{Title: "New", Count: 10, IsActive: false, Temperature: 25.3}

	resultB := graph.Replacer(currentB, changeB)

	if resultB.Title != "New" {
		t.Errorf("Expected Title='New', got '%s'", resultB.Title)
	}
	if resultB.Count != 10 {
		t.Errorf("Expected Count=10, got %d", resultB.Count)
	}
	if resultB.IsActive != false {
		t.Errorf("Expected IsActive=false, got %v", resultB.IsActive)
	}
	if resultB.Temperature != 25.3 {
		t.Errorf("Expected Temperature=25.3, got %f", resultB.Temperature)
	}
}

func TestReplacer_PointerFields(t *testing.T) {
	type StateWithPointer struct {
		Name  string
		Value *int
	}

	val1 := 42
	val2 := 99

	currentState := StateWithPointer{
		Name:  "current",
		Value: &val1,
	}

	changeState := StateWithPointer{
		Name:  "change",
		Value: &val2,
	}

	result := graph.Replacer(currentState, changeState)

	if result.Name != "change" {
		t.Errorf("Expected Name='change', got '%s'", result.Name)
	}

	if result.Value == nil {
		t.Fatal("Expected Value to be non-nil")
	}

	if *result.Value != 99 {
		t.Errorf("Expected *Value=99, got %d", *result.Value)
	}
}

func TestReplacer_NilPointerInChange(t *testing.T) {
	type StateWithPointer struct {
		Name  string
		Value *int
	}

	val1 := 42

	currentState := StateWithPointer{
		Name:  "current",
		Value: &val1,
	}

	changeState := StateWithPointer{
		Name:  "change",
		Value: nil,
	}

	result := graph.Replacer(currentState, changeState)

	if result.Name != "change" {
		t.Errorf("Expected Name='change', got '%s'", result.Name)
	}

	if result.Value != nil {
		t.Error("Expected Value to be nil after replacement")
	}
}

func TestReplacer_MultipleSequentialReplacements(t *testing.T) {
	state1 := SimpleState{Value: "first", Counter: 1}
	state2 := SimpleState{Value: "second", Counter: 2}
	state3 := SimpleState{Value: "third", Counter: 3}

	result1 := graph.Replacer(state1, state2)
	if result1.Value != "second" || result1.Counter != 2 {
		t.Errorf("First replacement failed: got %+v", result1)
	}

	result2 := graph.Replacer(result1, state3)
	if result2.Value != "third" || result2.Counter != 3 {
		t.Errorf("Second replacement failed: got %+v", result2)
	}

	// Original states should remain unchanged
	if state1.Value != "first" {
		t.Error("Original state1 was modified")
	}
}

func TestReplacer_SliceModification(t *testing.T) {
	type StateWithSlice struct {
		Items []string
	}

	currentState := StateWithSlice{
		Items: []string{"a", "b", "c"},
	}

	changeState := StateWithSlice{
		Items: []string{"x", "y"},
	}

	result := graph.Replacer(currentState, changeState)

	if len(result.Items) != 2 {
		t.Errorf("Expected Items length=2, got %d", len(result.Items))
	}

	if result.Items[0] != "x" || result.Items[1] != "y" {
		t.Errorf("Expected Items=['x', 'y'], got %v", result.Items)
	}
}

func TestReplacer_MapModification(t *testing.T) {
	type StateWithMap struct {
		Data map[string]int
	}

	currentState := StateWithMap{
		Data: map[string]int{"a": 1, "b": 2},
	}

	changeState := StateWithMap{
		Data: map[string]int{"c": 3},
	}

	result := graph.Replacer(currentState, changeState)

	if len(result.Data) != 1 {
		t.Errorf("Expected Data length=1, got %d", len(result.Data))
	}

	if val, ok := result.Data["c"]; !ok || val != 3 {
		t.Error("Expected Data to contain c=3")
	}

	if _, ok := result.Data["a"]; ok {
		t.Error("Expected Data to not contain old key 'a'")
	}
}
