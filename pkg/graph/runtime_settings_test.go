package graph_test

import (
	"testing"
	"time"

	"github.com/morphy76/ggraph/pkg/graph"
)

func TestFillRuntimeSettingsWithDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    graph.RuntimeSettings
		expected graph.RuntimeSettings
	}{
		{
			name:  "empty settings should use all defaults",
			input: graph.RuntimeSettings{},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom DefaultWorkerCount should override default",
			input: graph.RuntimeSettings{
				DefaultWorkerCount: 10,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             10,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom DefaultWorkerQueueSize should override default",
			input: graph.RuntimeSettings{
				DefaultWorkerQueueSize: 200,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         200,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom OutcomeNotificationQueueSize should override default",
			input: graph.RuntimeSettings{
				OutcomeNotificationQueueSize: 50,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   50,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom OutcomeNotificationMaxInterval should override default",
			input: graph.RuntimeSettings{
				OutcomeNotificationMaxInterval: 200 * time.Millisecond,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: 200 * time.Millisecond,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom PersistenceJobsQueueSize should override default",
			input: graph.RuntimeSettings{
				PersistenceJobsQueueSize: 20,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       20,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom PersistenceJobTimeout should override default",
			input: graph.RuntimeSettings{
				PersistenceJobTimeout: 10 * time.Second,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          10 * time.Second,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom ThreadTTL should override default",
			input: graph.RuntimeSettings{
				ThreadTTL: 2 * time.Hour,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      2 * time.Hour,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom ThreadEvictorInterval should override default",
			input: graph.RuntimeSettings{
				ThreadEvictorInterval: 10 * time.Minute,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          10 * time.Minute,
				GracefulShutdownTimeout:        graph.RuntimeSettingDefaultGracefulShutdownTimeout,
			},
		},
		{
			name: "custom GracefulShutdownTimeout should override default",
			input: graph.RuntimeSettings{
				GracefulShutdownTimeout: 30 * time.Second,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             graph.RuntimeSettingDefaultWorkerCount,
				DefaultWorkerQueueSize:         graph.RuntimeSettingDefaultWorkerQueueSize,
				OutcomeNotificationQueueSize:   graph.RuntimeSettingDefaultOutcomeNotificationQueueSize,
				OutcomeNotificationMaxInterval: graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval,
				PersistenceJobsQueueSize:       graph.RuntimeSettingDefaultPersistenceQueueSize,
				PersistenceJobTimeout:          graph.RuntimeSettingDefaultPersistenceTimeout,
				ThreadTTL:                      graph.RuntimeSettingDefaultThreadTTL,
				ThreadEvictorInterval:          graph.RuntimeSettingDefaultThreadEvictorInterval,
				GracefulShutdownTimeout:        30 * time.Second,
			},
		},
		{
			name: "all custom settings should override all defaults",
			input: graph.RuntimeSettings{
				DefaultWorkerCount:             20,
				DefaultWorkerQueueSize:         500,
				OutcomeNotificationQueueSize:   250,
				OutcomeNotificationMaxInterval: 500 * time.Millisecond,
				PersistenceJobsQueueSize:       50,
				PersistenceJobTimeout:          15 * time.Second,
				ThreadTTL:                      3 * time.Hour,
				ThreadEvictorInterval:          15 * time.Minute,
				GracefulShutdownTimeout:        60 * time.Second,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             20,
				DefaultWorkerQueueSize:         500,
				OutcomeNotificationQueueSize:   250,
				OutcomeNotificationMaxInterval: 500 * time.Millisecond,
				PersistenceJobsQueueSize:       50,
				PersistenceJobTimeout:          15 * time.Second,
				ThreadTTL:                      3 * time.Hour,
				ThreadEvictorInterval:          15 * time.Minute,
				GracefulShutdownTimeout:        60 * time.Second,
			},
		},
		{
			name: "negative values should be preserved",
			input: graph.RuntimeSettings{
				DefaultWorkerCount:             -1,
				DefaultWorkerQueueSize:         -1,
				OutcomeNotificationQueueSize:   -1,
				OutcomeNotificationMaxInterval: -1 * time.Second,
				PersistenceJobsQueueSize:       -1,
				PersistenceJobTimeout:          -1 * time.Second,
				ThreadTTL:                      -1 * time.Hour,
				ThreadEvictorInterval:          -1 * time.Minute,
				GracefulShutdownTimeout:        -1 * time.Second,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             -1,
				DefaultWorkerQueueSize:         -1,
				OutcomeNotificationQueueSize:   -1,
				OutcomeNotificationMaxInterval: -1 * time.Second,
				PersistenceJobsQueueSize:       -1,
				PersistenceJobTimeout:          -1 * time.Second,
				ThreadTTL:                      -1 * time.Hour,
				ThreadEvictorInterval:          -1 * time.Minute,
				GracefulShutdownTimeout:        -1 * time.Second,
			},
		},
		{
			name: "very small positive values should be preserved",
			input: graph.RuntimeSettings{
				DefaultWorkerCount:             1,
				DefaultWorkerQueueSize:         1,
				OutcomeNotificationQueueSize:   1,
				OutcomeNotificationMaxInterval: 1 * time.Nanosecond,
				PersistenceJobsQueueSize:       1,
				PersistenceJobTimeout:          1 * time.Nanosecond,
				ThreadTTL:                      1 * time.Nanosecond,
				ThreadEvictorInterval:          1 * time.Nanosecond,
				GracefulShutdownTimeout:        1 * time.Nanosecond,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             1,
				DefaultWorkerQueueSize:         1,
				OutcomeNotificationQueueSize:   1,
				OutcomeNotificationMaxInterval: 1 * time.Nanosecond,
				PersistenceJobsQueueSize:       1,
				PersistenceJobTimeout:          1 * time.Nanosecond,
				ThreadTTL:                      1 * time.Nanosecond,
				ThreadEvictorInterval:          1 * time.Nanosecond,
				GracefulShutdownTimeout:        1 * time.Nanosecond,
			},
		},
		{
			name: "very large values should be preserved",
			input: graph.RuntimeSettings{
				DefaultWorkerCount:             10000,
				DefaultWorkerQueueSize:         100000,
				OutcomeNotificationQueueSize:   50000,
				OutcomeNotificationMaxInterval: 24 * time.Hour,
				PersistenceJobsQueueSize:       5000,
				PersistenceJobTimeout:          1 * time.Hour,
				ThreadTTL:                      7 * 24 * time.Hour,
				ThreadEvictorInterval:          12 * time.Hour,
				GracefulShutdownTimeout:        5 * time.Minute,
			},
			expected: graph.RuntimeSettings{
				DefaultWorkerCount:             10000,
				DefaultWorkerQueueSize:         100000,
				OutcomeNotificationQueueSize:   50000,
				OutcomeNotificationMaxInterval: 24 * time.Hour,
				PersistenceJobsQueueSize:       5000,
				PersistenceJobTimeout:          1 * time.Hour,
				ThreadTTL:                      7 * 24 * time.Hour,
				ThreadEvictorInterval:          12 * time.Hour,
				GracefulShutdownTimeout:        5 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := graph.FillRuntimeSettingsWithDefaults(tt.input)

			if result.DefaultWorkerCount != tt.expected.DefaultWorkerCount {
				t.Errorf("DefaultWorkerCount = %d, want %d", result.DefaultWorkerCount, tt.expected.DefaultWorkerCount)
			}
			if result.DefaultWorkerQueueSize != tt.expected.DefaultWorkerQueueSize {
				t.Errorf("DefaultWorkerQueueSize = %d, want %d", result.DefaultWorkerQueueSize, tt.expected.DefaultWorkerQueueSize)
			}
			if result.OutcomeNotificationQueueSize != tt.expected.OutcomeNotificationQueueSize {
				t.Errorf("OutcomeNotificationQueueSize = %d, want %d", result.OutcomeNotificationQueueSize, tt.expected.OutcomeNotificationQueueSize)
			}
			if result.OutcomeNotificationMaxInterval != tt.expected.OutcomeNotificationMaxInterval {
				t.Errorf("OutcomeNotificationMaxInterval = %v, want %v", result.OutcomeNotificationMaxInterval, tt.expected.OutcomeNotificationMaxInterval)
			}
			if result.PersistenceJobsQueueSize != tt.expected.PersistenceJobsQueueSize {
				t.Errorf("PersistenceJobsQueueSize = %d, want %d", result.PersistenceJobsQueueSize, tt.expected.PersistenceJobsQueueSize)
			}
			if result.PersistenceJobTimeout != tt.expected.PersistenceJobTimeout {
				t.Errorf("PersistenceJobTimeout = %v, want %v", result.PersistenceJobTimeout, tt.expected.PersistenceJobTimeout)
			}
			if result.ThreadTTL != tt.expected.ThreadTTL {
				t.Errorf("ThreadTTL = %v, want %v", result.ThreadTTL, tt.expected.ThreadTTL)
			}
			if result.ThreadEvictorInterval != tt.expected.ThreadEvictorInterval {
				t.Errorf("ThreadEvictorInterval = %v, want %v", result.ThreadEvictorInterval, tt.expected.ThreadEvictorInterval)
			}
			if result.GracefulShutdownTimeout != tt.expected.GracefulShutdownTimeout {
				t.Errorf("GracefulShutdownTimeout = %v, want %v", result.GracefulShutdownTimeout, tt.expected.GracefulShutdownTimeout)
			}
		})
	}
}

func TestFillRuntimeSettingsWithDefaults_Idempotency(t *testing.T) {
	// Test that applying defaults multiple times produces the same result
	input := graph.RuntimeSettings{
		DefaultWorkerCount:      15,
		ThreadTTL:               2 * time.Hour,
		GracefulShutdownTimeout: 20 * time.Second,
	}

	first := graph.FillRuntimeSettingsWithDefaults(input)
	second := graph.FillRuntimeSettingsWithDefaults(first)

	if first.DefaultWorkerCount != second.DefaultWorkerCount {
		t.Errorf("Idempotency check failed for DefaultWorkerCount: first=%d, second=%d", first.DefaultWorkerCount, second.DefaultWorkerCount)
	}
	if first.DefaultWorkerQueueSize != second.DefaultWorkerQueueSize {
		t.Errorf("Idempotency check failed for DefaultWorkerQueueSize: first=%d, second=%d", first.DefaultWorkerQueueSize, second.DefaultWorkerQueueSize)
	}
	if first.OutcomeNotificationQueueSize != second.OutcomeNotificationQueueSize {
		t.Errorf("Idempotency check failed for OutcomeNotificationQueueSize: first=%d, second=%d", first.OutcomeNotificationQueueSize, second.OutcomeNotificationQueueSize)
	}
	if first.OutcomeNotificationMaxInterval != second.OutcomeNotificationMaxInterval {
		t.Errorf("Idempotency check failed for OutcomeNotificationMaxInterval: first=%v, second=%v", first.OutcomeNotificationMaxInterval, second.OutcomeNotificationMaxInterval)
	}
	if first.PersistenceJobsQueueSize != second.PersistenceJobsQueueSize {
		t.Errorf("Idempotency check failed for PersistenceJobsQueueSize: first=%d, second=%d", first.PersistenceJobsQueueSize, second.PersistenceJobsQueueSize)
	}
	if first.PersistenceJobTimeout != second.PersistenceJobTimeout {
		t.Errorf("Idempotency check failed for PersistenceJobTimeout: first=%v, second=%v", first.PersistenceJobTimeout, second.PersistenceJobTimeout)
	}
	if first.ThreadTTL != second.ThreadTTL {
		t.Errorf("Idempotency check failed for ThreadTTL: first=%v, second=%v", first.ThreadTTL, second.ThreadTTL)
	}
	if first.ThreadEvictorInterval != second.ThreadEvictorInterval {
		t.Errorf("Idempotency check failed for ThreadEvictorInterval: first=%v, second=%v", first.ThreadEvictorInterval, second.ThreadEvictorInterval)
	}
	if first.GracefulShutdownTimeout != second.GracefulShutdownTimeout {
		t.Errorf("Idempotency check failed for GracefulShutdownTimeout: first=%v, second=%v", first.GracefulShutdownTimeout, second.GracefulShutdownTimeout)
	}
}

func TestFillRuntimeSettingsWithDefaults_DoesNotMutateInput(t *testing.T) {
	// Test that the input settings are not modified
	input := graph.RuntimeSettings{
		DefaultWorkerCount:             0,
		DefaultWorkerQueueSize:         0,
		OutcomeNotificationQueueSize:   0,
		OutcomeNotificationMaxInterval: 0,
		PersistenceJobsQueueSize:       0,
		PersistenceJobTimeout:          0,
		ThreadTTL:                      0,
		ThreadEvictorInterval:          0,
		GracefulShutdownTimeout:        0,
	}

	original := input

	_ = graph.FillRuntimeSettingsWithDefaults(input)

	if input.DefaultWorkerCount != original.DefaultWorkerCount {
		t.Errorf("Input DefaultWorkerCount was mutated")
	}
	if input.DefaultWorkerQueueSize != original.DefaultWorkerQueueSize {
		t.Errorf("Input DefaultWorkerQueueSize was mutated")
	}
	if input.OutcomeNotificationQueueSize != original.OutcomeNotificationQueueSize {
		t.Errorf("Input OutcomeNotificationQueueSize was mutated")
	}
	if input.OutcomeNotificationMaxInterval != original.OutcomeNotificationMaxInterval {
		t.Errorf("Input OutcomeNotificationMaxInterval was mutated")
	}
	if input.PersistenceJobsQueueSize != original.PersistenceJobsQueueSize {
		t.Errorf("Input PersistenceJobsQueueSize was mutated")
	}
	if input.PersistenceJobTimeout != original.PersistenceJobTimeout {
		t.Errorf("Input PersistenceJobTimeout was mutated")
	}
	if input.ThreadTTL != original.ThreadTTL {
		t.Errorf("Input ThreadTTL was mutated")
	}
	if input.ThreadEvictorInterval != original.ThreadEvictorInterval {
		t.Errorf("Input ThreadEvictorInterval was mutated")
	}
	if input.GracefulShutdownTimeout != original.GracefulShutdownTimeout {
		t.Errorf("Input GracefulShutdownTimeout was mutated")
	}
}

func TestRuntimeSettingDefaultConstants(t *testing.T) {
	// Test that the default constants have expected values
	if graph.RuntimeSettingDefaultWorkerCount <= 0 {
		t.Errorf("RuntimeSettingDefaultWorkerCount should be positive, got %d", graph.RuntimeSettingDefaultWorkerCount)
	}
	if graph.RuntimeSettingDefaultWorkerQueueSize <= 0 {
		t.Errorf("RuntimeSettingDefaultWorkerQueueSize should be positive, got %d", graph.RuntimeSettingDefaultWorkerQueueSize)
	}
	if graph.RuntimeSettingDefaultOutcomeNotificationQueueSize <= 0 {
		t.Errorf("RuntimeSettingDefaultOutcomeNotificationQueueSize should be positive, got %d", graph.RuntimeSettingDefaultOutcomeNotificationQueueSize)
	}
	if graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval <= 0 {
		t.Errorf("RuntimeSettingDefaultOutcomeNotificationMaxInterval should be positive, got %v", graph.RuntimeSettingDefaultOutcomeNotificationMaxInterval)
	}
	if graph.RuntimeSettingDefaultPersistenceQueueSize <= 0 {
		t.Errorf("RuntimeSettingDefaultPersistenceQueueSize should be positive, got %d", graph.RuntimeSettingDefaultPersistenceQueueSize)
	}
	if graph.RuntimeSettingDefaultPersistenceTimeout <= 0 {
		t.Errorf("RuntimeSettingDefaultPersistenceTimeout should be positive, got %v", graph.RuntimeSettingDefaultPersistenceTimeout)
	}
	if graph.RuntimeSettingDefaultThreadTTL <= 0 {
		t.Errorf("RuntimeSettingDefaultThreadTTL should be positive, got %v", graph.RuntimeSettingDefaultThreadTTL)
	}
	if graph.RuntimeSettingDefaultThreadEvictorInterval <= 0 {
		t.Errorf("RuntimeSettingDefaultThreadEvictorInterval should be positive, got %v", graph.RuntimeSettingDefaultThreadEvictorInterval)
	}
	if graph.RuntimeSettingDefaultGracefulShutdownTimeout <= 0 {
		t.Errorf("RuntimeSettingDefaultGracefulShutdownTimeout should be positive, got %v", graph.RuntimeSettingDefaultGracefulShutdownTimeout)
	}
}
