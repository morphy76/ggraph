package agent

import (
	"testing"
	"time"
)

func TestMessageRole(t *testing.T) {
	tests := []struct {
		name string
		role MessageRole
		want MessageRole
	}{
		{
			name: "system role",
			role: System,
			want: 0,
		},
		{
			name: "user role",
			role: User,
			want: 1,
		},
		{
			name: "assistant role",
			role: Assistant,
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.role != tt.want {
				t.Errorf("MessageRole = %v, want %v", tt.role, tt.want)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		message Message
		verify  func(*testing.T, Message)
	}{
		{
			name: "system message with timestamp",
			message: Message{
				Ts:      now,
				Role:    System,
				Content: "You are a helpful assistant.",
			},
			verify: func(t *testing.T, msg Message) {
				if msg.Role != System {
					t.Errorf("Expected System role, got %v", msg.Role)
				}
				if msg.Content != "You are a helpful assistant." {
					t.Errorf("Expected specific content, got '%s'", msg.Content)
				}
				if !msg.Ts.Equal(now) {
					t.Errorf("Expected timestamp %v, got %v", now, msg.Ts)
				}
			},
		},
		{
			name: "user message",
			message: Message{
				Ts:      now,
				Role:    User,
				Content: "Hello!",
			},
			verify: func(t *testing.T, msg Message) {
				if msg.Role != User {
					t.Errorf("Expected User role, got %v", msg.Role)
				}
				if msg.Content != "Hello!" {
					t.Errorf("Expected 'Hello!', got '%s'", msg.Content)
				}
			},
		},
		{
			name: "assistant message",
			message: Message{
				Ts:      now,
				Role:    Assistant,
				Content: "How can I help you?",
			},
			verify: func(t *testing.T, msg Message) {
				if msg.Role != Assistant {
					t.Errorf("Expected Assistant role, got %v", msg.Role)
				}
				if msg.Content != "How can I help you?" {
					t.Errorf("Expected specific content, got '%s'", msg.Content)
				}
			},
		},
		{
			name: "empty content",
			message: Message{
				Ts:      now,
				Role:    User,
				Content: "",
			},
			verify: func(t *testing.T, msg Message) {
				if msg.Content != "" {
					t.Errorf("Expected empty content, got '%s'", msg.Content)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.verify != nil {
				tt.verify(t, tt.message)
			}
		})
	}
}

func TestConversation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		conversation Conversation
		verify       func(*testing.T, Conversation)
	}{
		{
			name: "empty conversation",
			conversation: Conversation{
				Messages: []Message{},
			},
			verify: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 0 {
					t.Errorf("Expected 0 messages, got %d", len(conv.Messages))
				}
			},
		},
		{
			name: "conversation with single message",
			conversation: Conversation{
				Messages: []Message{
					{Ts: now, Role: System, Content: "You are a helpful assistant."},
				},
			},
			verify: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 1 {
					t.Errorf("Expected 1 message, got %d", len(conv.Messages))
				}
				if conv.Messages[0].Role != System {
					t.Errorf("Expected System role, got %v", conv.Messages[0].Role)
				}
			},
		},
		{
			name: "conversation with multiple messages",
			conversation: Conversation{
				Messages: []Message{
					{Ts: now, Role: System, Content: "You are a helpful assistant."},
					{Ts: now.Add(time.Second), Role: User, Content: "Hello!"},
					{Ts: now.Add(2 * time.Second), Role: Assistant, Content: "Hi there!"},
				},
			},
			verify: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 3 {
					t.Errorf("Expected 3 messages, got %d", len(conv.Messages))
				}
				expectedRoles := []MessageRole{System, User, Assistant}
				for i, expectedRole := range expectedRoles {
					if conv.Messages[i].Role != expectedRole {
						t.Errorf("Message %d: Expected role %v, got %v", i, expectedRole, conv.Messages[i].Role)
					}
				}
				// Verify chronological order
				for i := 1; i < len(conv.Messages); i++ {
					if conv.Messages[i].Ts.Before(conv.Messages[i-1].Ts) {
						t.Errorf("Messages are not in chronological order at index %d", i)
					}
				}
			},
		},
		{
			name: "conversation with multiple user messages",
			conversation: Conversation{
				Messages: []Message{
					{Ts: now, Role: User, Content: "First question"},
					{Ts: now.Add(time.Second), Role: User, Content: "Second question"},
					{Ts: now.Add(2 * time.Second), Role: User, Content: "Third question"},
				},
			},
			verify: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 3 {
					t.Errorf("Expected 3 messages, got %d", len(conv.Messages))
				}
				for i, msg := range conv.Messages {
					if msg.Role != User {
						t.Errorf("Message %d: Expected User role, got %v", i, msg.Role)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.verify != nil {
				tt.verify(t, tt.conversation)
			}
		})
	}
}

func TestConversationAppend(t *testing.T) {
	// Test appending messages to a conversation
	conv := Conversation{
		Messages: []Message{},
	}

	msg1 := Message{Ts: time.Now(), Role: System, Content: "System message"}
	conv.Messages = append(conv.Messages, msg1)

	if len(conv.Messages) != 1 {
		t.Errorf("Expected 1 message after first append, got %d", len(conv.Messages))
	}

	msg2 := Message{Ts: time.Now(), Role: User, Content: "User message"}
	conv.Messages = append(conv.Messages, msg2)

	if len(conv.Messages) != 2 {
		t.Errorf("Expected 2 messages after second append, got %d", len(conv.Messages))
	}

	if conv.Messages[0].Role != System {
		t.Errorf("Expected first message to be System, got %v", conv.Messages[0].Role)
	}
	if conv.Messages[1].Role != User {
		t.Errorf("Expected second message to be User, got %v", conv.Messages[1].Role)
	}
}
