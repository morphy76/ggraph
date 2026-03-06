# Summarization Example

This directory contains an example demonstrating how to use the conversation summarization feature in `ggraph`.

## Overview

When building conversational agents, the context window can quickly grow large. This is an issue for:
- **Cost**: APIs often charge per token. Sending long histories costs more.
- **Context Limits**: Models have maximum token limits.
- **Latency**: Processing large contexts takes longer.

To solve this, `ggraph` includes a built-in `SummarizationConfig` that can automatically shrink the conversation history by summarizing older messages into a single system prompt while preserving recent turns for immediate context.

## Running the Example

In this directory, run:

```bash
export AIW_API_KEY="your-api-key"
go run run.go
```

The example sets up an interactive loop with a very low threshold (`MessageThreshold: 3`, `KeepRecentCount: 1`). As you interact with the agent, notice how the `[Current Conversation Messages Length]` value resets after reaching the threshold. The older messages are compressed into a single `System` message containing the summary of the conversation.

## Customization

The summarization feature is fully customizable:

```go
summaryConfig := &a.SummarizationConfig{
    Enabled:          true,
    MessageThreshold: 20, // Summarize after 20 messages
    KeepRecentCount:  5, // Keep the last 5 messages unsummarized
    PromptOverride:   "Custom summarizer prompt...",
    Summarizer:       nil, // Set a custom `SummarizerFn` logic if needed
}

chatNode, err := o.CreateConversationNode(
    "ChatNode",
    "velvet-14b",
    aiwClient,
    a.WithSummarization(summaryConfig), // Apply the config
)
```
