package agent

import (
	"testing"
)

func TestWithBestOf(t *testing.T) {
	tests := []struct {
		name    string
		bestOf  int64
		wantErr bool
	}{
		{
			name:    "valid best of",
			bestOf:  3,
			wantErr: false,
		},
		{
			name:    "minimum valid value",
			bestOf:  1,
			wantErr: false,
		},
		{
			name:    "invalid - zero",
			bestOf:  0,
			wantErr: true,
		},
		{
			name:    "invalid - negative",
			bestOf:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBestOf(tt.bestOf)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithBestOf().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.BestOf == nil || *opts.BestOf != tt.bestOf) {
				if opts.BestOf == nil {
					t.Errorf("Expected BestOf %d, got nil", tt.bestOf)
				} else {
					t.Errorf("Expected BestOf %d, got %d", tt.bestOf, *opts.BestOf)
				}
			}
		})
	}
}

func TestWithEcho(t *testing.T) {
	tests := []struct {
		name string
		echo bool
	}{
		{
			name: "echo true",
			echo: true,
		},
		{
			name: "echo false",
			echo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEcho(tt.echo)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if err != nil {
				t.Errorf("WithEcho().ApplyToCompletion() unexpected error = %v", err)
			}
			if opts.Echo == nil || *opts.Echo != tt.echo {
				if opts.Echo == nil {
					t.Errorf("Expected Echo %v, got nil", tt.echo)
				} else {
					t.Errorf("Expected Echo %v, got %v", tt.echo, *opts.Echo)
				}
			}
		})
	}
}

func TestWithFrequencyPenalty(t *testing.T) {
	tests := []struct {
		name             string
		frequencyPenalty float64
		wantErr          bool
	}{
		{
			name:             "valid positive value",
			frequencyPenalty: 1.0,
			wantErr:          false,
		},
		{
			name:             "valid negative value",
			frequencyPenalty: -1.5,
			wantErr:          false,
		},
		{
			name:             "valid zero",
			frequencyPenalty: 0.0,
			wantErr:          false,
		},
		{
			name:             "valid maximum",
			frequencyPenalty: 2.0,
			wantErr:          false,
		},
		{
			name:             "valid minimum",
			frequencyPenalty: -2.0,
			wantErr:          false,
		},
		{
			name:             "invalid - too high",
			frequencyPenalty: 2.1,
			wantErr:          true,
		},
		{
			name:             "invalid - too low",
			frequencyPenalty: -2.1,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithFrequencyPenalty(tt.frequencyPenalty)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithFrequencyPenalty().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.FrequencyPenalty == nil || *opts.FrequencyPenalty != tt.frequencyPenalty) {
				if opts.FrequencyPenalty == nil {
					t.Errorf("Expected FrequencyPenalty %f, got nil", tt.frequencyPenalty)
				} else {
					t.Errorf("Expected FrequencyPenalty %f, got %f", tt.frequencyPenalty, *opts.FrequencyPenalty)
				}
			}
		})
	}
}

func TestWithLogprobs(t *testing.T) {
	tests := []struct {
		name     string
		logprobs int64
		wantErr  bool
	}{
		{
			name:     "valid value",
			logprobs: 3,
			wantErr:  false,
		},
		{
			name:     "minimum valid",
			logprobs: 0,
			wantErr:  false,
		},
		{
			name:     "maximum valid",
			logprobs: 5,
			wantErr:  false,
		},
		{
			name:     "invalid - negative",
			logprobs: -1,
			wantErr:  true,
		},
		{
			name:     "invalid - too high",
			logprobs: 6,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLogprobs(tt.logprobs)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithLogprobs().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.Logprobs == nil || *opts.Logprobs != tt.logprobs) {
				if opts.Logprobs == nil {
					t.Errorf("Expected Logprobs %d, got nil", tt.logprobs)
				} else {
					t.Errorf("Expected Logprobs %d, got %d", tt.logprobs, *opts.Logprobs)
				}
			}
		})
	}
}

func TestWithMaxTokens(t *testing.T) {
	tests := []struct {
		name      string
		maxTokens int64
		wantErr   bool
	}{
		{
			name:      "valid value",
			maxTokens: 100,
			wantErr:   false,
		},
		{
			name:      "minimum valid",
			maxTokens: 1,
			wantErr:   false,
		},
		{
			name:      "large valid value",
			maxTokens: 10000,
			wantErr:   false,
		},
		{
			name:      "invalid - zero",
			maxTokens: 0,
			wantErr:   true,
		},
		{
			name:      "invalid - negative",
			maxTokens: -1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMaxTokens(tt.maxTokens)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithMaxTokens().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.MaxTokens == nil || *opts.MaxTokens != tt.maxTokens) {
				if opts.MaxTokens == nil {
					t.Errorf("Expected MaxTokens %d, got nil", tt.maxTokens)
				} else {
					t.Errorf("Expected MaxTokens %d, got %d", tt.maxTokens, *opts.MaxTokens)
				}
			}
		})
	}
}

func TestWithN(t *testing.T) {
	tests := []struct {
		name    string
		n       int64
		wantErr bool
	}{
		{
			name:    "valid value",
			n:       2,
			wantErr: false,
		},
		{
			name:    "minimum valid",
			n:       1,
			wantErr: false,
		},
		{
			name:    "large valid value",
			n:       10,
			wantErr: false,
		},
		{
			name:    "invalid - zero",
			n:       0,
			wantErr: true,
		},
		{
			name:    "invalid - negative",
			n:       -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithN(tt.n)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithN().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.N == nil || *opts.N != tt.n) {
				if opts.N == nil {
					t.Errorf("Expected N %d, got nil", tt.n)
				} else {
					t.Errorf("Expected N %d, got %d", tt.n, *opts.N)
				}
			}
		})
	}
}

func TestWithPresencePenalty(t *testing.T) {
	tests := []struct {
		name            string
		presencePenalty float64
		wantErr         bool
	}{
		{
			name:            "valid positive value",
			presencePenalty: 1.0,
			wantErr:         false,
		},
		{
			name:            "valid negative value",
			presencePenalty: -1.5,
			wantErr:         false,
		},
		{
			name:            "valid zero",
			presencePenalty: 0.0,
			wantErr:         false,
		},
		{
			name:            "valid maximum",
			presencePenalty: 2.0,
			wantErr:         false,
		},
		{
			name:            "valid minimum",
			presencePenalty: -2.0,
			wantErr:         false,
		},
		{
			name:            "invalid - too high",
			presencePenalty: 2.1,
			wantErr:         true,
		},
		{
			name:            "invalid - too low",
			presencePenalty: -2.1,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPresencePenalty(tt.presencePenalty)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithPresencePenalty().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.PresencePenalty == nil || *opts.PresencePenalty != tt.presencePenalty) {
				if opts.PresencePenalty == nil {
					t.Errorf("Expected PresencePenalty %f, got nil", tt.presencePenalty)
				} else {
					t.Errorf("Expected PresencePenalty %f, got %f", tt.presencePenalty, *opts.PresencePenalty)
				}
			}
		})
	}
}

func TestWithSeed(t *testing.T) {
	tests := []struct {
		name string
		seed int64
	}{
		{
			name: "positive value",
			seed: 42,
		},
		{
			name: "zero",
			seed: 0,
		},
		{
			name: "negative value",
			seed: -100,
		},
		{
			name: "large value",
			seed: 999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithSeed(tt.seed)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if err != nil {
				t.Errorf("WithSeed().ApplyToCompletion() unexpected error = %v", err)
			}
			if opts.Seed == nil || *opts.Seed != tt.seed {
				if opts.Seed == nil {
					t.Errorf("Expected Seed %d, got nil", tt.seed)
				} else {
					t.Errorf("Expected Seed %d, got %d", tt.seed, *opts.Seed)
				}
			}
		})
	}
}

func TestWithTemperature(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		wantErr     bool
	}{
		{
			name:        "valid low value",
			temperature: 0.2,
			wantErr:     false,
		},
		{
			name:        "valid mid value",
			temperature: 0.7,
			wantErr:     false,
		},
		{
			name:        "valid high value",
			temperature: 1.8,
			wantErr:     false,
		},
		{
			name:        "valid minimum",
			temperature: 0.0,
			wantErr:     false,
		},
		{
			name:        "valid maximum",
			temperature: 2.0,
			wantErr:     false,
		},
		{
			name:        "invalid - negative",
			temperature: -0.1,
			wantErr:     true,
		},
		{
			name:        "invalid - too high",
			temperature: 2.1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTemperature(tt.temperature)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithTemperature().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.Temperature == nil || *opts.Temperature != tt.temperature) {
				if opts.Temperature == nil {
					t.Errorf("Expected Temperature %f, got nil", tt.temperature)
				} else {
					t.Errorf("Expected Temperature %f, got %f", tt.temperature, *opts.Temperature)
				}
			}
		})
	}
}

func TestWithTopP(t *testing.T) {
	tests := []struct {
		name    string
		topP    float64
		wantErr bool
	}{
		{
			name:    "valid low value",
			topP:    0.1,
			wantErr: false,
		},
		{
			name:    "valid mid value",
			topP:    0.5,
			wantErr: false,
		},
		{
			name:    "valid high value",
			topP:    0.9,
			wantErr: false,
		},
		{
			name:    "valid minimum",
			topP:    0.0,
			wantErr: false,
		},
		{
			name:    "valid maximum",
			topP:    1.0,
			wantErr: false,
		},
		{
			name:    "invalid - negative",
			topP:    -0.1,
			wantErr: true,
		},
		{
			name:    "invalid - too high",
			topP:    1.1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTopP(tt.topP)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithTopP().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.TopP == nil || *opts.TopP != tt.topP) {
				if opts.TopP == nil {
					t.Errorf("Expected TopP %f, got nil", tt.topP)
				} else {
					t.Errorf("Expected TopP %f, got %f", tt.topP, *opts.TopP)
				}
			}
		})
	}
}

func TestWithUser(t *testing.T) {
	tests := []struct {
		name string
		user string
	}{
		{
			name: "normal user id",
			user: "user-1234",
		},
		{
			name: "empty user id",
			user: "",
		},
		{
			name: "complex user id",
			user: "user-abc-123-xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithUser(tt.user)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if err != nil {
				t.Errorf("WithUser().ApplyToCompletion() unexpected error = %v", err)
			}
			if opts.User == nil || *opts.User != tt.user {
				if opts.User == nil {
					t.Errorf("Expected User '%s', got nil", tt.user)
				} else {
					t.Errorf("Expected User '%s', got '%s'", tt.user, *opts.User)
				}
			}
		})
	}
}

func TestWithMaxCompletionTokens(t *testing.T) {
	tests := []struct {
		name                string
		maxCompletionTokens int64
		wantErr             bool
	}{
		{
			name:                "valid value",
			maxCompletionTokens: 200,
			wantErr:             false,
		},
		{
			name:                "minimum valid",
			maxCompletionTokens: 1,
			wantErr:             false,
		},
		{
			name:                "large valid value",
			maxCompletionTokens: 5000,
			wantErr:             false,
		},
		{
			name:                "invalid - zero",
			maxCompletionTokens: 0,
			wantErr:             true,
		},
		{
			name:                "invalid - negative",
			maxCompletionTokens: -1,
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMaxCompletionTokens(tt.maxCompletionTokens)
			opts := &ModelOptions{}
			err := opt.ApplyToCompletion(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("WithMaxCompletionTokens().ApplyToCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (opts.MaxCompletionTokens == nil || *opts.MaxCompletionTokens != tt.maxCompletionTokens) {
				if opts.MaxCompletionTokens == nil {
					t.Errorf("Expected MaxCompletionTokens %d, got nil", tt.maxCompletionTokens)
				} else {
					t.Errorf("Expected MaxCompletionTokens %d, got %d", tt.maxCompletionTokens, *opts.MaxCompletionTokens)
				}
			}
		})
	}
}

func TestModelOptionApplyToConversation(t *testing.T) {
	tests := []struct {
		name     string
		option   ModelOption
		wantErr  bool
		validate func(*testing.T, *ModelOptions)
	}{
		{
			name:    "temperature applies to conversation",
			option:  WithTemperature(0.8),
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.Temperature == nil || *opts.Temperature != 0.8 {
					if opts.Temperature == nil {
						t.Errorf("Expected Temperature 0.8, got nil")
					} else {
						t.Errorf("Expected Temperature 0.8, got %f", *opts.Temperature)
					}
				}
			},
		},
		{
			name:    "invalid temperature fails on conversation",
			option:  WithTemperature(3.0),
			wantErr: true,
		},
		{
			name:    "max completion tokens applies to conversation",
			option:  WithMaxCompletionTokens(300),
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.MaxCompletionTokens == nil || *opts.MaxCompletionTokens != 300 {
					if opts.MaxCompletionTokens == nil {
						t.Errorf("Expected MaxCompletionTokens 300, got nil")
					} else {
						t.Errorf("Expected MaxCompletionTokens 300, got %d", *opts.MaxCompletionTokens)
					}
				}
			},
		},
		{
			name:    "presence penalty applies to conversation",
			option:  WithPresencePenalty(0.5),
			wantErr: false,
			validate: func(t *testing.T, opts *ModelOptions) {
				if opts.PresencePenalty == nil || *opts.PresencePenalty != 0.5 {
					if opts.PresencePenalty == nil {
						t.Errorf("Expected PresencePenalty 0.5, got nil")
					} else {
						t.Errorf("Expected PresencePenalty 0.5, got %f", *opts.PresencePenalty)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ModelOptions{}
			err := tt.option.ApplyToConversation(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyToConversation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, opts)
			}
		})
	}
}

func TestModelOptionChaining(t *testing.T) {
	// Test that multiple options can be applied sequentially
	opts := &ModelOptions{}

	options := []ModelOption{
		WithTemperature(0.7),
		WithMaxTokens(150),
		WithTopP(0.9),
		WithN(2),
		WithFrequencyPenalty(0.5),
		WithPresencePenalty(-0.5),
		WithSeed(42),
		WithUser("test-user"),
	}

	for _, opt := range options {
		if err := opt.ApplyToCompletion(opts); err != nil {
			t.Errorf("Failed to apply option: %v", err)
		}
	}

	// Verify all options were applied
	if opts.Temperature == nil || *opts.Temperature != 0.7 {
		if opts.Temperature == nil {
			t.Errorf("Expected Temperature 0.7, got nil")
		} else {
			t.Errorf("Expected Temperature 0.7, got %f", *opts.Temperature)
		}
	}
	if opts.MaxTokens == nil || *opts.MaxTokens != 150 {
		if opts.MaxTokens == nil {
			t.Errorf("Expected MaxTokens 150, got nil")
		} else {
			t.Errorf("Expected MaxTokens 150, got %d", *opts.MaxTokens)
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
	if opts.FrequencyPenalty == nil || *opts.FrequencyPenalty != 0.5 {
		if opts.FrequencyPenalty == nil {
			t.Errorf("Expected FrequencyPenalty 0.5, got nil")
		} else {
			t.Errorf("Expected FrequencyPenalty 0.5, got %f", *opts.FrequencyPenalty)
		}
	}
	if opts.PresencePenalty == nil || *opts.PresencePenalty != -0.5 {
		if opts.PresencePenalty == nil {
			t.Errorf("Expected PresencePenalty -0.5, got nil")
		} else {
			t.Errorf("Expected PresencePenalty -0.5, got %f", *opts.PresencePenalty)
		}
	}
	if opts.Seed == nil || *opts.Seed != 42 {
		if opts.Seed == nil {
			t.Errorf("Expected Seed 42, got nil")
		} else {
			t.Errorf("Expected Seed 42, got %d", *opts.Seed)
		}
	}
	if opts.User == nil || *opts.User != "test-user" {
		if opts.User == nil {
			t.Errorf("Expected User 'test-user', got nil")
		} else {
			t.Errorf("Expected User 'test-user', got '%s'", *opts.User)
		}
	}
}
