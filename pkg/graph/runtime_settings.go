package graph

import "time"

const (
	// RuntimeSettingDefaultWorkerCount is the default number of workers in the runtime.
	RuntimeSettingDefaultWorkerCount = 5
	// RuntimeSettingDefaultWorkerQueueSize is the default size of the worker queue in the runtime.
	RuntimeSettingDefaultWorkerQueueSize = 100

	// RuntimeSettingDefaultOutcomeNotificationQueueSize is the default size of the outcome notification queue used to communicate outside the graph.
	RuntimeSettingDefaultOutcomeNotificationQueueSize = 100
	// RuntimeSettingDefaultOutcomeNotificationMaxInterval is the default maximum interval between outcome notifications.
	RuntimeSettingDefaultOutcomeNotificationMaxInterval = 100 * time.Millisecond

	// RuntimeSettingDefaultPersistenceQueueSize is the default size of the queue in the runtime worker which flushes pending states.
	RuntimeSettingDefaultPersistenceQueueSize = 10
	// RuntimeSettingDefaultPersistenceTimeout is the default timeout between persistence flushes.
	RuntimeSettingDefaultPersistenceTimeout = 5 * time.Second

	// RuntimeSettingDefaultThreadTTL is the default time-to-live for inactive threads.
	RuntimeSettingDefaultThreadTTL = 1 * time.Hour
	// RuntimeSettingDefaultThreadEvictorInterval is the default interval for evicting inactive threads.
	RuntimeSettingDefaultThreadEvictorInterval = 5 * time.Minute

	// RuntimeSettingDefaultGracefulShutdownTimeout is the default timeout for graceful shutdown operations.
	RuntimeSettingDefaultGracefulShutdownTimeout = 10 * time.Second
)

// RuntimeSettings holds the configuration settings for the graph runtime.
type RuntimeSettings struct {
	// DefaultWorkerCount is the default number of workers in the runtime.
	DefaultWorkerCount int
	// DefaultWorkerQueueSize is the default size of the worker queue in the runtime.
	DefaultWorkerQueueSize int

	// OutcomeNotificationQueueSize is the default size of the outcome notification queue used to communicate outside the graph.
	OutcomeNotificationQueueSize int
	// OutcomeNotificationMaxInterval is the default maximum interval between outcome notifications.
	OutcomeNotificationMaxInterval time.Duration

	// PersistenceJobsQueueSize is the default size of the queue in the runtime worker which flushes pending states.
	PersistenceJobsQueueSize int
	// PersistenceJobTimeout is the default timeout between persistence flushes.
	PersistenceJobTimeout time.Duration

	// ThreadTTL is the default time-to-live for inactive threads.
	ThreadTTL time.Duration
	// ThreadEvictorInterval is the default interval for evicting inactive threads.
	ThreadEvictorInterval time.Duration

	// GracefulShutdownTimeout is the default timeout for graceful shutdown operations.
	GracefulShutdownTimeout time.Duration
}

var defaultRuntimeSettings = RuntimeSettings{
	DefaultWorkerCount:     RuntimeSettingDefaultWorkerCount,
	DefaultWorkerQueueSize: RuntimeSettingDefaultWorkerQueueSize,

	OutcomeNotificationQueueSize:   RuntimeSettingDefaultOutcomeNotificationQueueSize,
	OutcomeNotificationMaxInterval: RuntimeSettingDefaultOutcomeNotificationMaxInterval,

	PersistenceJobsQueueSize: RuntimeSettingDefaultPersistenceQueueSize,
	PersistenceJobTimeout:    RuntimeSettingDefaultPersistenceTimeout,

	ThreadTTL:             RuntimeSettingDefaultThreadTTL,
	ThreadEvictorInterval: RuntimeSettingDefaultThreadEvictorInterval,

	GracefulShutdownTimeout: RuntimeSettingDefaultGracefulShutdownTimeout,
}

// FillRuntimeSettingsWithDefaults fills in any zero-value settings with their default values.
func FillRuntimeSettingsWithDefaults(s RuntimeSettings) RuntimeSettings {
	merged := defaultRuntimeSettings

	if s.DefaultWorkerCount != 0 {
		merged.DefaultWorkerCount = s.DefaultWorkerCount
	}
	if s.DefaultWorkerQueueSize != 0 {
		merged.DefaultWorkerQueueSize = s.DefaultWorkerQueueSize
	}

	if s.OutcomeNotificationQueueSize != 0 {
		merged.OutcomeNotificationQueueSize = s.OutcomeNotificationQueueSize
	}
	if s.OutcomeNotificationMaxInterval != 0 {
		merged.OutcomeNotificationMaxInterval = s.OutcomeNotificationMaxInterval
	}

	if s.PersistenceJobsQueueSize != 0 {
		merged.PersistenceJobsQueueSize = s.PersistenceJobsQueueSize
	}
	if s.PersistenceJobTimeout != 0 {
		merged.PersistenceJobTimeout = s.PersistenceJobTimeout
	}

	if s.ThreadTTL != 0 {
		merged.ThreadTTL = s.ThreadTTL
	}
	if s.ThreadEvictorInterval != 0 {
		merged.ThreadEvictorInterval = s.ThreadEvictorInterval
	}

	if s.GracefulShutdownTimeout != 0 {
		merged.GracefulShutdownTimeout = s.GracefulShutdownTimeout
	}

	return merged
}
