package agent

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultSummarizer(t *testing.T) {
	now := time.Now()

	t.Run("below threshold does nothing", func(t *testing.T) {
		messages := []Message{
			{Ts: now, Role: System, Content: "Sys msg"},
			{Ts: now, Role: User, Content: "User msg 1"},
			{Ts: now, Role: Assistant, Content: "Asst msg 1"},
		}

		result, err := DefaultSummarizer(messages, 5, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 3 {
			t.Fatalf("expected 3 messages, got %d", len(result))
		}
	})

	t.Run("above threshold summarizes", func(t *testing.T) {
		messages := []Message{
			{Ts: now, Role: System, Content: "Sys msg"},
			{Ts: now, Role: User, Content: "Q1"},
			{Ts: now, Role: Assistant, Content: "A1"},
			{Ts: now, Role: User, Content: "Q2"},
			{Ts: now, Role: Assistant, Content: "A2"},
			{Ts: now, Role: User, Content: "Q3"},
			{Ts: now, Role: Assistant, Content: "A3"}, // 7 messages total
		}

		result, err := DefaultSummarizer(messages, 5, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Expected format:
		// 1 System message with summary
		// Plus keep 2 recent messages
		// Total: 3 messages
		if len(result) != 3 {
			t.Fatalf("expected 3 messages, got %d", len(result))
		}

		if result[0].Role != System {
			t.Errorf("expected first message to be System, got %v", result[0].Role)
		}

		if !strings.Contains(result[0].Content, "Previous conversation summary:") {
			t.Errorf("expected summary prefix, got %s", result[0].Content)
		}

		if !strings.Contains(result[0].Content, "Summarized 5 messages (2 user, 2 assistant, 0 tool, 1 system)") {
			t.Errorf("expected counts, got %s", result[0].Content)
		}

		if !strings.Contains(result[0].Content, "Initial user request: Q1") {
			t.Errorf("expected initial user request, got %s", result[0].Content)
		}

		if !strings.Contains(result[0].Content, "Last assistant response: A2") {
			t.Errorf("expected last assistant response, got %s", result[0].Content)
		}

		// Last 2 string should be kept unchanged
		if result[1].Content != "Q3" || result[2].Content != "A3" {
			t.Errorf("expected latest 2 original messages kept, got %s and %s", result[1].Content, result[2].Content)
		}
	})

	t.Run("empty conversation", func(t *testing.T) {
		messages := []Message{}

		result, err := DefaultSummarizer(messages, 5, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 0 {
			t.Fatalf("expected 0 messages, got %d", len(result))
		}
	})

	t.Run("fewer messages than KeepRecentCount", func(t *testing.T) {
		messages := []Message{
			{Ts: now, Role: User, Content: "Q1"},
		}

		result, err := DefaultSummarizer(messages, 0, 5) // threshold 0 to force run but recentCount > len
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("expected 1 messages, got %d", len(result))
		}
	})
}
