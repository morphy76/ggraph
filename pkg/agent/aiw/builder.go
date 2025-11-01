package aiw

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	o "github.com/morphy76/ggraph/pkg/agent/openai"
)

// NewAIWClient creates a new OpenAI client configured for the AIW platform.
//
// Parameters:
//   - PAT: The Personal Access Token (PAT) for authentication.
//   - opts: Additional request options for the OpenAI API calls.
//
// Returns:
//   - A pointer to an instance of openai.Client configured for AIW.
//
// Example usage:
//
//	client := NewAIWClient("your-api-key", option.WithTimeout(30*time.Second))
func NewAIWClient(
	PAT string,
	opts ...option.RequestOption,
) *openai.Client {
	return o.NewClient(AIWBaseURL, PAT, opts...)
}
