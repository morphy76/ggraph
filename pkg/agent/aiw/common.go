package aiw

import "os"

const (
	// AIWBaseURL is the base URL for the Almawave AIW Platform.
	AIWBaseURL = "https://portal.aiwave.ai/llm/api"
	// EnvKeyPAT is the environment variable key for the AIW API key.
	EnvKeyPAT = "AIW_API_KEY"
)

// PATFromEnv retrieves the AIW API key from the environment variable "AIW_API_KEY
//
// Returns:
//   - The Personal Access Token (PAT) as a string.
func PATFromEnv() string {
	return os.Getenv(EnvKeyPAT)
}
