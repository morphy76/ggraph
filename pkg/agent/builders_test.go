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

func TestCreateCompletionOptions(t *testing.T) {
	tests := []struct {
		name              string
		prompt            string
		model             string
		completionOptions []ModelOption
		wantErr           bool
		validate          func(*testing.T, *ModelOptions)
	}{
		{
			name:              "basic options without extra parameters",
			prompt:            "Test prompt",
			model:             "gpt-4",
			completionOptions: []ModelOption{},
			wantErr:           false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Prompt != "Test prompt" {
					t.Errorf("Expected prompt 'Test prompt', got '%s'", opts.Prompt)
				}
				if opts.Model != "gpt-4" {
					t.Errorf("Expected model 'gpt-4', got '%s'", opts.Model)
				}
			},
		},
		{
			name:   "with valid max tokens",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithMaxTokens(100),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.MaxTokens == nil || *opts.MaxTokens != 100 {
					if opts.MaxTokens == nil {
						t.Errorf("Expected MaxTokens 100, got nil")
					} else {
						t.Errorf("Expected MaxTokens 100, got %d", *opts.MaxTokens)
					}
				}
			},
		},
		{
			name:   "with invalid max tokens",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithMaxTokens(0),
			},
			wantErr: true,
		},
		{
			name:   "with valid temperature",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithTemperature(0.7),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Temperature == nil || *opts.Temperature != 0.7 {
					if opts.Temperature == nil {
						t.Errorf("Expected Temperature 0.7, got nil")
					} else {
						t.Errorf("Expected Temperature 0.7, got %f", *opts.Temperature)
					}
				}
			},
		},
		{
			name:   "with invalid temperature",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithTemperature(3.0),
			},
			wantErr: true,
		},
		{
			name:   "with multiple valid options",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithMaxTokens(150),
				WithTemperature(0.8),
				WithTopP(0.9),
				WithN(2),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.MaxTokens == nil || *opts.MaxTokens != 150 {
					if opts.MaxTokens == nil {
						t.Errorf("Expected MaxTokens 150, got nil")
					} else {
						t.Errorf("Expected MaxTokens 150, got %d", *opts.MaxTokens)
					}
				}
				if opts.Temperature == nil || *opts.Temperature != 0.8 {
					if opts.Temperature == nil {
						t.Errorf("Expected Temperature 0.8, got nil")
					} else {
						t.Errorf("Expected Temperature 0.8, got %f", *opts.Temperature)
					}
				}
				if opts.TopP == nil || *opts.TopP != 0.9 {
					if opts.TopP == nil {
						t.Errorf("Expected TopP 0.9, got nil")
					} else {
						t.Errorf("Expected TopP 0.9, got %f", *opts.TopP)
					}
				}
				if opts.N == nil || *opts.N != 2 {
					if opts.N == nil {
						t.Errorf("Expected N 2, got nil")
					} else {
						t.Errorf("Expected N 2, got %d", *opts.N)
					}
				}
			},
		},
		{
			name:   "with one invalid option among valid ones",
			prompt: "Test prompt",
			model:  "gpt-4",
			completionOptions: []ModelOption{
				WithMaxTokens(150),
				WithTemperature(-1.0), // Invalid
				WithTopP(0.9),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateCompletionOptions(tt.model, tt.prompt, tt.completionOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCompletionOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, got)
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

func TestCreateConversationOptions(t *testing.T) {
	testMessages := []Message{
		CreateMessage(System, "You are a helpful assistant."),
		CreateMessage(User, "Tell me a joke."),
	}

	tests := []struct {
		name         string
		promptModel  string
		messages     []Message
		modelOptions []ModelOption
		wantErr      bool
		validate     func(*testing.T, *ModelOptions)
	}{
		{
			name:         "basic options without extra parameters",
			promptModel:  "gpt-4-chat",
			messages:     testMessages,
			modelOptions: []ModelOption{},
			wantErr:      false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Model != "gpt-4-chat" {
					t.Errorf("Expected model 'gpt-4-chat', got '%s'", opts.Model)
				}
			},
		},
		{
			name:        "with valid temperature",
			promptModel: "gpt-4-chat",
			messages:    testMessages,
			modelOptions: []ModelOption{
				WithTemperature(0.7),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Temperature == nil || *opts.Temperature != 0.7 {
					if opts.Temperature == nil {
						t.Errorf("Expected Temperature 0.7, got nil")
					} else {
						t.Errorf("Expected Temperature 0.7, got %f", *opts.Temperature)
					}
				}
			},
		},
		{
			name:        "with invalid temperature",
			promptModel: "gpt-4-chat",
			messages:    testMessages,
			modelOptions: []ModelOption{
				WithTemperature(2.5),
			},
			wantErr: true,
		},
		{
			name:        "with valid max completion tokens",
			promptModel: "gpt-4-chat",
			messages:    testMessages,
			modelOptions: []ModelOption{
				WithMaxCompletionTokens(200),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.MaxCompletionTokens == nil || *opts.MaxCompletionTokens != 200 {
					if opts.MaxCompletionTokens == nil {
						t.Errorf("Expected MaxCompletionTokens 200, got nil")
					} else {
						t.Errorf("Expected MaxCompletionTokens 200, got %d", *opts.MaxCompletionTokens)
					}
				}
			},
		},
		{
			name:        "with multiple valid options",
			promptModel: "gpt-4-chat",
			messages:    testMessages,
			modelOptions: []ModelOption{
				WithTemperature(0.8),
				WithMaxCompletionTokens(300),
				WithPresencePenalty(0.5),
				WithFrequencyPenalty(-0.5),
			},
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Temperature == nil || *opts.Temperature != 0.8 {
					if opts.Temperature == nil {
						t.Errorf("Expected Temperature 0.8, got nil")
					} else {
						t.Errorf("Expected Temperature 0.8, got %f", *opts.Temperature)
					}
				}
				if opts.MaxCompletionTokens == nil || *opts.MaxCompletionTokens != 300 {
					if opts.MaxCompletionTokens == nil {
						t.Errorf("Expected MaxCompletionTokens 300, got nil")
					} else {
						t.Errorf("Expected MaxCompletionTokens 300, got %d", *opts.MaxCompletionTokens)
					}
				}
				if opts.PresencePenalty == nil || *opts.PresencePenalty != 0.5 {
					if opts.PresencePenalty == nil {
						t.Errorf("Expected PresencePenalty 0.5, got nil")
					} else {
						t.Errorf("Expected PresencePenalty 0.5, got %f", *opts.PresencePenalty)
					}
				}
				if opts.FrequencyPenalty == nil || *opts.FrequencyPenalty != -0.5 {
					if opts.FrequencyPenalty == nil {
						t.Errorf("Expected FrequencyPenalty -0.5, got nil")
					} else {
						t.Errorf("Expected FrequencyPenalty -0.5, got %f", *opts.FrequencyPenalty)
					}
				}
			},
		},
		{
			name:        "with one invalid option among valid ones",
			promptModel: "gpt-4-chat",
			messages:    testMessages,
			modelOptions: []ModelOption{
				WithTemperature(0.8),
				WithPresencePenalty(3.0), // Invalid
			},
			wantErr: true,
		},
		{
			name:         "empty messages",
			promptModel:  "gpt-4-chat",
			messages:     []Message{},
			modelOptions: []ModelOption{},
			wantErr:      false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Model != "gpt-4-chat" {
					t.Errorf("Expected model 'gpt-4-chat', got '%s'", opts.Model)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateConversationOptions(tt.promptModel, tt.messages, tt.modelOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConversationOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}
