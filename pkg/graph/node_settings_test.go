package graph_test

import (
	"testing"
	"time"

	"github.com/morphy76/ggraph/pkg/graph"
)

func TestFillNodeSettingsWithDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    graph.NodeSettings
		expected graph.NodeSettings
	}{
		{
			name:  "empty settings should use all defaults",
			input: graph.NodeSettings{},
			expected: graph.NodeSettings{
				MailboxSize:   graph.NodeSettingDefaultMailboxSize,
				AcceptTimeout: graph.NodeSettingDefaultAcceptTimeout,
			},
		},
		{
			name: "custom MailboxSize should override default",
			input: graph.NodeSettings{
				MailboxSize: 50,
			},
			expected: graph.NodeSettings{
				MailboxSize:   50,
				AcceptTimeout: graph.NodeSettingDefaultAcceptTimeout,
			},
		},
		{
			name: "custom AcceptTimeout should override default",
			input: graph.NodeSettings{
				AcceptTimeout: 10 * time.Second,
			},
			expected: graph.NodeSettings{
				MailboxSize:   graph.NodeSettingDefaultMailboxSize,
				AcceptTimeout: 10 * time.Second,
			},
		},
		{
			name: "all custom settings should override all defaults",
			input: graph.NodeSettings{
				MailboxSize:   100,
				AcceptTimeout: 15 * time.Second,
			},
			expected: graph.NodeSettings{
				MailboxSize:   100,
				AcceptTimeout: 15 * time.Second,
			},
		},
		{
			name: "negative MailboxSize should use default",
			input: graph.NodeSettings{
				MailboxSize:   -1,
				AcceptTimeout: 3 * time.Second,
			},
			expected: graph.NodeSettings{
				MailboxSize:   -1,
				AcceptTimeout: 3 * time.Second,
			},
		},
		{
			name: "negative AcceptTimeout should use value as-is",
			input: graph.NodeSettings{
				MailboxSize:   20,
				AcceptTimeout: -1 * time.Second,
			},
			expected: graph.NodeSettings{
				MailboxSize:   20,
				AcceptTimeout: -1 * time.Second,
			},
		},
		{
			name: "very small MailboxSize should be preserved",
			input: graph.NodeSettings{
				MailboxSize: 1,
			},
			expected: graph.NodeSettings{
				MailboxSize:   1,
				AcceptTimeout: graph.NodeSettingDefaultAcceptTimeout,
			},
		},
		{
			name: "very large MailboxSize should be preserved",
			input: graph.NodeSettings{
				MailboxSize: 10000,
			},
			expected: graph.NodeSettings{
				MailboxSize:   10000,
				AcceptTimeout: graph.NodeSettingDefaultAcceptTimeout,
			},
		},
		{
			name: "very small AcceptTimeout should be preserved",
			input: graph.NodeSettings{
				AcceptTimeout: 1 * time.Nanosecond,
			},
			expected: graph.NodeSettings{
				MailboxSize:   graph.NodeSettingDefaultMailboxSize,
				AcceptTimeout: 1 * time.Nanosecond,
			},
		},
		{
			name: "very large AcceptTimeout should be preserved",
			input: graph.NodeSettings{
				AcceptTimeout: 24 * time.Hour,
			},
			expected: graph.NodeSettings{
				MailboxSize:   graph.NodeSettingDefaultMailboxSize,
				AcceptTimeout: 24 * time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := graph.FillNodeSettingsWithDefaults(tt.input)

			if result.MailboxSize != tt.expected.MailboxSize {
				t.Errorf("MailboxSize = %d, want %d", result.MailboxSize, tt.expected.MailboxSize)
			}

			if result.AcceptTimeout != tt.expected.AcceptTimeout {
				t.Errorf("AcceptTimeout = %v, want %v", result.AcceptTimeout, tt.expected.AcceptTimeout)
			}
		})
	}
}

func TestFillNodeSettingsWithDefaults_Idempotency(t *testing.T) {
	// Test that applying defaults multiple times produces the same result
	input := graph.NodeSettings{
		MailboxSize: 25,
	}

	first := graph.FillNodeSettingsWithDefaults(input)
	second := graph.FillNodeSettingsWithDefaults(first)

	if first.MailboxSize != second.MailboxSize {
		t.Errorf("Idempotency check failed for MailboxSize: first=%d, second=%d", first.MailboxSize, second.MailboxSize)
	}

	if first.AcceptTimeout != second.AcceptTimeout {
		t.Errorf("Idempotency check failed for AcceptTimeout: first=%v, second=%v", first.AcceptTimeout, second.AcceptTimeout)
	}
}

func TestFillNodeSettingsWithDefaults_DoesNotMutateInput(t *testing.T) {
	// Test that the input settings are not modified
	input := graph.NodeSettings{
		MailboxSize:   0,
		AcceptTimeout: 0,
	}

	originalMailboxSize := input.MailboxSize
	originalAcceptTimeout := input.AcceptTimeout

	_ = graph.FillNodeSettingsWithDefaults(input)

	if input.MailboxSize != originalMailboxSize {
		t.Errorf("Input MailboxSize was mutated: was=%d, now=%d", originalMailboxSize, input.MailboxSize)
	}

	if input.AcceptTimeout != originalAcceptTimeout {
		t.Errorf("Input AcceptTimeout was mutated: was=%v, now=%v", originalAcceptTimeout, input.AcceptTimeout)
	}
}

func TestNodeSettingDefaultConstants(t *testing.T) {
	// Test that the default constants have expected values
	if graph.NodeSettingDefaultMailboxSize <= 0 {
		t.Errorf("NodeSettingDefaultMailboxSize should be positive, got %d", graph.NodeSettingDefaultMailboxSize)
	}

	if graph.NodeSettingDefaultAcceptTimeout <= 0 {
		t.Errorf("NodeSettingDefaultAcceptTimeout should be positive, got %v", graph.NodeSettingDefaultAcceptTimeout)
	}
}
