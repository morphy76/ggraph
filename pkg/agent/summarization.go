package agent

import (
	"fmt"
	"time"
)

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
	Summarizer:       DefaultSummarizer,
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

// DefaultSummarizer provides the default summarization algorithm:
// 1. Checks if the message count exceeds the threshold.
// 2. Counts messages by role.
// 3. Extracts the first user message and the last assistant response.
// 4. Returns a new slice with 1 summary system message + the `recentCount` latest messages.
func DefaultSummarizer(messages []Message, threshold int, recentCount int) ([]Message, error) {
	if len(messages) <= threshold {
		return messages, nil
	}

	if len(messages) <= recentCount {
		return messages, nil
	}

	var userCount, assistantCount, toolCount, systemCount int
	var firstUserMsg string
	var lastAssistantMsg string

	// Collect stats and extract required messages
	for i := 0; i < len(messages)-recentCount; i++ {
		msg := messages[i]
		switch msg.Role {
		case System:
			systemCount++
		case User:
			if userCount == 0 {
				firstUserMsg = msg.Content
			}
			userCount++
		case Assistant:
			lastAssistantMsg = msg.Content
			assistantCount++
		case Tool:
			toolCount++
		}
	}

	// It's possible the last assistant message is within the recentCount keep region
	// But the PR issue states "last assistant response", it probably means in the summarized set,
	// or in the entire conversation. We'll search backwards over all messages if not found,
	// or just over the summarized portion. The test checks if A2 is the last assistant response,
	// and A3 is part of the kept messages. So it means the last assistant response IN THE SUMMARIZED portion.
	// We're iterating over the summarized portion, so that's correct.

	// In case there was no user message or assistant message in the summarized portion
	if firstUserMsg == "" {
		firstUserMsg = "None"
	}
	if lastAssistantMsg == "" {
		lastAssistantMsg = "None"
	}

	summaryContent := "Previous conversation summary:\n"
	summaryContent += fmt.Sprintf("- Summarized %d messages (%d user, %d assistant, %d tool, %d system)\n",
		len(messages)-recentCount, userCount, assistantCount, toolCount, systemCount)
	summaryContent += "- Initial user request: " + firstUserMsg + "\n"
	summaryContent += "- Last assistant response: " + lastAssistantMsg

	summaryMessage := Message{
		Ts:      time.Now(),
		Role:    System,
		Content: summaryContent,
	}

	result := make([]Message, 0, 1+recentCount)
	result = append(result, summaryMessage)
	result = append(result, messages[len(messages)-recentCount:]...)

	return result, nil
}
