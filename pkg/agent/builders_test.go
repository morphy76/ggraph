package agent

import (
	"testing"
	"time"
)

func TestCreateCompletion(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Completion
	}{
		{
			name: "simple text",
			text: "This is a test completion",
			want: Completion{Text: "This is a test completion"},
		},
		{
			name: "empty text",
			text: "",
			want: Completion{Text: ""},
		},
		{
			name: "multiline text",
			text: "Line 1\nLine 2\nLine 3",
			want: Completion{Text: "Line 1\nLine 2\nLine 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateCompletion(tt.text)
			if got.Text != tt.want.Text {
				t.Errorf("CreateCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	beforeTime := time.Now()
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name    string
		role    MessageRole
		content string
	}{
		{
			name:    "system message",
			role:    System,
			content: "You are a helpful assistant.",
		},
		{
			name:    "user message",
			role:    User,
			content: "What's the weather like?",
		},
		{
			name:    "assistant message",
			role:    Assistant,
			content: "I can help you with that.",
		},
		{
			name:    "empty content",
			role:    User,
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMessage(tt.role, tt.content)

			if got.Role != tt.role {
				t.Errorf("CreateMessage() Role = %v, want %v", got.Role, tt.role)
			}
			if got.Content != tt.content {
				t.Errorf("CreateMessage() Content = %v, want %v", got.Content, tt.content)
			}
			if got.Ts.Before(beforeTime) {
				t.Errorf("CreateMessage() Ts should be after test start time")
			}
		})
	}
}

func TestCreateConversation(t *testing.T) {
	tests := []struct {
		name     string
		messages []Message
		validate func(*testing.T, Conversation)
	}{
		{
			name:     "empty conversation",
			messages: []Message{},
			validate: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 0 {
					t.Errorf("Expected empty messages, got %d messages", len(conv.Messages))
				}
			},
		},
		{
			name: "single message",
			messages: []Message{
				CreateMessage(System, "You are a helpful assistant."),
			},
			validate: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 1 {
					t.Errorf("Expected 1 message, got %d messages", len(conv.Messages))
				}
				if conv.Messages[0].Role != System {
					t.Errorf("Expected System role, got %v", conv.Messages[0].Role)
				}
				if conv.Messages[0].Content != "You are a helpful assistant." {
					t.Errorf("Expected specific content, got '%s'", conv.Messages[0].Content)
				}
			},
		},
		{
			name: "multiple messages",
			messages: []Message{
				CreateMessage(System, "You are a helpful assistant."),
				CreateMessage(User, "Hello!"),
				CreateMessage(Assistant, "Hi there! How can I help you?"),
			},
			validate: func(t *testing.T, conv Conversation) {
				if len(conv.Messages) != 3 {
					t.Errorf("Expected 3 messages, got %d messages", len(conv.Messages))
				}
				expectedRoles := []MessageRole{System, User, Assistant}
				for i, expectedRole := range expectedRoles {
					if conv.Messages[i].Role != expectedRole {
						t.Errorf("Message %d: Expected role %v, got %v", i, expectedRole, conv.Messages[i].Role)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateConversation(tt.messages...)
			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}
