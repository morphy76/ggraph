package agent

import (
	"testing"
)

func TestCompletion(t *testing.T) {
	tests := []struct {
		name       string
		completion Completion
		verify     func(*testing.T, Completion)
	}{
		{
			name: "simple completion",
			completion: Completion{
				Text: "This is a generated response.",
			},
			verify: func(t *testing.T, comp Completion) {
				if comp.Text != "This is a generated response." {
					t.Errorf("Expected specific text, got '%s'", comp.Text)
				}
			},
		},
		{
			name: "empty completion",
			completion: Completion{
				Text: "",
			},
			verify: func(t *testing.T, comp Completion) {
				if comp.Text != "" {
					t.Errorf("Expected empty text, got '%s'", comp.Text)
				}
			},
		},
		{
			name: "multiline completion",
			completion: Completion{
				Text: "Line 1\nLine 2\nLine 3",
			},
			verify: func(t *testing.T, comp Completion) {
				expected := "Line 1\nLine 2\nLine 3"
				if comp.Text != expected {
					t.Errorf("Expected multiline text, got '%s'", comp.Text)
				}
			},
		},
		{
			name: "completion with special characters",
			completion: Completion{
				Text: "Hello, 世界! こんにちは 🌍",
			},
			verify: func(t *testing.T, comp Completion) {
				expected := "Hello, 世界! こんにちは 🌍"
				if comp.Text != expected {
					t.Errorf("Expected text with special characters, got '%s'", comp.Text)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.verify != nil {
				tt.verify(t, tt.completion)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are defined
	tests := []struct {
		name  string
		err   error
		check func(error) bool
	}{
		{
			name: "ErrorInvalidBestOf is defined",
			err:  ErrorInvalidBestOf,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidFrequencyPenalty is defined",
			err:  ErrorInvalidFrequencyPenalty,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidLogprobs is defined",
			err:  ErrorInvalidLogprobs,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidMaxTokens is defined",
			err:  ErrorInvalidMaxTokens,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidN is defined",
			err:  ErrorInvalidN,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidPresencePenalty is defined",
			err:  ErrorInvalidPresencePenalty,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidTemperature is defined",
			err:  ErrorInvalidTemperature,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
		{
			name: "ErrorInvalidTopP is defined",
			err:  ErrorInvalidTopP,
			check: func(err error) bool {
				return err != nil && err.Error() != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(tt.err) {
				t.Errorf("Error constant check failed for %s", tt.name)
			}
		})
	}
}

func TestErrorUniqueness(t *testing.T) {
	// Test that all error messages are unique
	errors := []error{
		ErrorInvalidBestOf,
		ErrorInvalidFrequencyPenalty,
		ErrorInvalidLogprobs,
		ErrorInvalidMaxTokens,
		ErrorInvalidN,
		ErrorInvalidPresencePenalty,
		ErrorInvalidTemperature,
		ErrorInvalidTopP,
	}

	errorMessages := make(map[string]bool)
	for _, err := range errors {
		msg := err.Error()
		if errorMessages[msg] {
			t.Errorf("Duplicate error message found: %s", msg)
		}
		errorMessages[msg] = true
	}
}
