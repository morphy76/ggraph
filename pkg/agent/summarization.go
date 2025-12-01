package agent

// SummarizerFn defines a function that summarizes a sequence of messages
type SummarizerFn func(messages []Message, threshold int, recentCount int) ([]Message, error)

// SummarizationConfig holds the configuration for conversation summarization
type SummarizationConfig struct {
	// Whether summarization is active
	Enabled bool
	// Optional prompt override for the summarization model
	PromptOverride string
	// Number of messages that triggers summarization
	MessageThreshold int
	// Number of recent messages to keep unsummarized
	KeepRecentCount int
	// Custom summarizer function (nil uses default)
	Summarizer SummarizerFn
}

var defaultSummarizationConfig = &SummarizationConfig{
	MessageThreshold: 20,
	PromptOverride:   "Summarize this conversation between a user and an AI assistant, preserving important details and context. Feel free to drop trivial or courtesy related exchanges. Keep the language. Be concise yet comprehensive.",
	KeepRecentCount:  2,
	Summarizer:       nil,
}

// FillSummarizationConfigWithDefaults fills in default values for any zero-value fields in the provided SummarizationConfig.
//
// Parameters:
//   - cfg: A pointer to a SummarizationConfig instance.
//
// Returns:
//   - A pointer to a SummarizationConfig instance with defaults filled in.
//
// Example usage:
//
//	cfg := &SummarizationConfig{Enabled: true}
//	filledCfg := FillSummarizationConfigWithDefaults(cfg)
func FillSummarizationConfigWithDefaults(cfg *SummarizationConfig) *SummarizationConfig {
	if cfg == nil {
		return defaultSummarizationConfig
	}
	if cfg.PromptOverride == "" {
		cfg.PromptOverride = defaultSummarizationConfig.PromptOverride
	}
	if cfg.MessageThreshold == 0 {
		cfg.MessageThreshold = defaultSummarizationConfig.MessageThreshold
	}
	if cfg.KeepRecentCount == 0 {
		cfg.KeepRecentCount = defaultSummarizationConfig.KeepRecentCount
	}
	if cfg.Summarizer == nil {
		cfg.Summarizer = defaultSummarizationConfig.Summarizer
	}
	return cfg
}
