# LLM Package TODOs

**Package:** `pkg/llm` and subpackages (`ollama/`, `openai/`, `aiw/`)  
**Date:** October 15, 2025  
**Status:** Functional but needs enhancements

---

## 🔴 Critical Priority

### Core Message Model

- [ ] **Multimodal Content Support** (High Impact)
  - Current: Only `string` content supported
  - Need: Support for images, files, audio, video
  - Suggested API:
    ```go
    type ContentType string
    const (
        TextContent ContentType = "text"
        ImageContent ContentType = "image"
        AudioContent ContentType = "audio"
        VideoContent ContentType = "video"
        FileContent ContentType = "file"
    )
    
    type MessageContent interface {
        Type() ContentType
        Data() interface{}
    }
    
    type TextMessageContent struct { Text string }
    type ImageMessageContent struct { URL string; Data []byte; MimeType string }
    
    type Message struct {
        Ts       time.Time
        Role     MessageRole
        Content  []MessageContent  // Changed from string to slice
        Metadata map[string]string
    }
    ```

- [ ] **Tool/Function Calling Support**
  - OpenAI and Ollama support function/tool calling
  - Add `ToolCall` and `ToolResult` message types
  - Example:
    ```go
    type ToolCall struct {
        ID       string
        Name     string
        Args     map[string]interface{}
    }
    
    type Message struct {
        // ... existing fields
        ToolCalls []ToolCall
        ToolCallID string  // For tool responses
    }
    ```

### Error Handling

- [ ] **Streaming Error Handling**
  - File: `ollama/ollama_chat.go:35`
  - Current: `respFunc` never returns error
  - Problem: Streaming errors are silently ignored
  - Fix: Properly handle and propagate streaming errors

- [ ] **Context Propagation**
  - All chat functions use `context.Background()`
  - Should accept context from graph runtime
  - Enables proper cancellation and timeout control

---

## 🟡 High Priority

### Configuration & Options

- [ ] **Configurable Timeouts**
  - Current: Hard-coded 60s timeout
  - Need: Per-provider, per-call configuration
  - Suggested:
    ```go
    type ChatOptions struct {
        Timeout     time.Duration
        Temperature float64
        MaxTokens   int
        TopP        float64
        Stream      bool
    }
    
    func CreateOllamaChatNode(name, model string, opts ChatOptions) (g.Node[llm.AgentModel], error)
    ```

- [ ] **Provider Configuration**
  - Current: Environment variables only (`OLLAMA_HOST`, `OPENAI_API_KEY`)
  - Need: Explicit configuration option
  - Suggested:
    ```go
    type OllamaConfig struct {
        Host    string
        Timeout time.Duration
    }
    
    func CreateOllamaChatNodeWithConfig(name, model string, config OllamaConfig) (g.Node[llm.AgentModel], error)
    ```

- [ ] **AIW Base URL Configuration** 🔥
  - File: `aiw/aiw_chat.go:18`
  - Current: `os.Setenv("OPENAI_BASE_URL", ...)` - **BREAKS OpenAI CLIENT**
  - Problem: Environment variable pollution affects all OpenAI clients
  - Impact: Cannot use multiple providers simultaneously
  - Fix: Use OpenAI client options directly:
    ```go
    client := openai.NewClient(
        openai.WithBaseURL("https://portal.aiwave.ai/llm/api"),
        openai.WithAPIKey(apiKey),
    )
    ```

### Node Interface Flexibility

- [ ] **Improve Node Interfaces for Better Adoption**
  - Current: Single function signature for all node types
  - Need: More flexible node creation patterns
  - Suggested approaches:
    1. **Options Pattern:**
       ```go
       type NodeOption func(*NodeConfig)
       
       func WithTemperature(temp float64) NodeOption
       func WithMaxTokens(tokens int) NodeOption
       func WithSystemPrompt(prompt string) NodeOption
       
       node := ollama.CreateChatNode("chat", "llama2",
           WithTemperature(0.7),
           WithMaxTokens(500),
       )
       ```
    
    2. **Builder Pattern:**
       ```go
       node := ollama.NewChatNodeBuilder("chat").
           Model("llama2").
           Temperature(0.7).
           MaxTokens(500).
           SystemPrompt("You are helpful").
           Build()
       ```
    
    3. **Config Struct:**
       ```go
       config := ollama.ChatConfig{
           Model:        "llama2",
           Temperature:  0.7,
           MaxTokens:    500,
           SystemPrompt: "You are helpful",
       }
       node := ollama.CreateChatNodeWithConfig("chat", config)
       ```

### Message Management

- [ ] **Add Helper Methods to AgentModel**
  - Current: Only `AddUserMessage(content string)`
  - Need:
    ```go
    func (m *AgentModel) AddSystemMessage(content string)
    func (m *AgentModel) AddAssistantMessage(content string)
    func (m *AgentModel) GetLastMessage() *Message
    func (m *AgentModel) GetMessagesByRole(role MessageRole) []Message
    func (m *AgentModel) Clear()
    func (m *AgentModel) Truncate(maxMessages int)
    func (m *AgentModel) TokenCount() (int, error)  // Estimate token count
    ```

- [ ] **Message Validation**
  - Validate message sequences (e.g., no consecutive assistant messages)
  - Validate content length
  - Validate required fields

### Partial Updates (Streaming)

- [ ] **Clarify Partial Update Semantics**
  - File: `ollama/ollama_chat.go:37-40`, `openai/openai_chat.go:29-34`
  - Current: `notify()` receives single message, unclear how runtime handles it
  - Problem: No clear contract for partial vs. full updates
  - Need: Document expected behavior or create typed partial updates

- [ ] **Streaming Progress Indicators**
  - Add token count/estimate in partial updates
  - Add timing information
  - Add stop reason when stream completes

---

## 🟢 Medium Priority

### Testing

- [ ] **Unit Tests for Message Conversion**
  - Test `ToLLamaModel`, `ToOpenAIModel`, etc.
  - Test role encoding/decoding
  - Test edge cases (empty messages, unknown roles)

- [ ] **Integration Tests**
  - Mock LLM responses for testing
  - Test streaming behavior
  - Test error scenarios

- [ ] **Benchmark Tests**
  - Message conversion performance
  - Memory allocations in streaming
  - Large conversation handling

### Provider-Specific Enhancements

#### Ollama (`pkg/llm/ollama/`)

- [ ] **Support Additional Ollama Options**
  - `temperature`, `top_p`, `top_k`, `repeat_penalty`
  - `num_ctx` (context window size)
  - `num_predict` (max tokens to generate)
  - Format options (json mode)

- [ ] **Embedding Support**
  - Create `CreateOllamaEmbeddingNode` for RAG patterns
  - Support batch embeddings

- [ ] **Model Management**
  - List available models
  - Pull/download models
  - Check model existence

#### OpenAI (`pkg/llm/openai/`)

- [ ] **Support Additional OpenAI Features**
  - Response format (JSON mode)
  - Logprobs
  - Multiple choices (n > 1)
  - Presence/frequency penalties
  - Stop sequences

- [ ] **Vision Support** (GPT-4 Vision)
  - Image input support
  - Image URL vs. base64 handling

- [ ] **Function/Tool Calling**
  - Define tool schemas
  - Handle tool calls in responses
  - Execute tools and continue conversation

- [ ] **Embeddings Support**
  - text-embedding-ada-002 and newer models

#### AIW (`pkg/llm/aiw/`)

- [ ] **Add Streaming Support**
  - Currently uses synchronous `New()` instead of `NewStreaming()`
  - Add streaming version for consistency with other providers

- [ ] **Document AIW-Specific Features**
  - Available models
  - Rate limits
  - Differences from standard OpenAI API

### Code Quality

- [ ] **Add Package Documentation**
  - Create `doc.go` for `pkg/llm`
  - Document message model design
  - Provide usage examples

- [ ] **Consistent Error Messages**
  - Define standard error types
  - Include provider context in errors
  - Wrap upstream errors consistently

- [ ] **Export Important Constants**
  - Default timeouts
  - Max message lengths
  - Token limits per model

---

## 🔵 Low Priority / Nice to Have

### Advanced Features

- [ ] **Conversation Management**
  - Automatic conversation summarization when too long
  - Token counting and management
  - Sliding window for long conversations

- [ ] **Prompt Templates**
  - Template engine for system prompts
  - Variable substitution
  - Example library of common prompts

- [ ] **Response Caching**
  - Cache identical requests
  - Configurable cache TTL
  - Cache invalidation strategies

- [ ] **Retry Logic**
  - Exponential backoff for rate limits
  - Automatic retry on transient errors
  - Circuit breaker pattern

- [ ] **Cost Tracking**
  - Token usage tracking per message
  - Cost estimation per provider/model
  - Budget limits

- [ ] **Multi-Provider Support**
  - Abstract provider interface
  - Provider selection strategy (fallback, load balancing)
  - Unified configuration

### Observability

- [ ] **Metrics**
  - Request count, latency, errors per provider
  - Token usage metrics
  - Streaming vs. non-streaming split

- [ ] **Logging**
  - Structured logging of requests/responses
  - Configurable log levels
  - Sensitive data redaction

- [ ] **Tracing**
  - OpenTelemetry integration
  - Span per LLM call
  - Trace conversation context

### Additional Provider Support

- [ ] **Anthropic Claude**
  - Create `pkg/llm/anthropic/` package
  - Support Claude-specific features (prompt caching, extended context)

- [ ] **Google Gemini**
  - Create `pkg/llm/gemini/` package
  - Multimodal support

- [ ] **Cohere**
  - Create `pkg/llm/cohere/` package
  - RAG-specific features

- [ ] **Azure OpenAI**
  - Separate from standard OpenAI (different auth, endpoints)
  - Deployment-based model selection

- [ ] **Local Models**
  - llama.cpp integration
  - GGUF model loading
  - GPU acceleration support

---

## 📋 Refactoring Ideas

### 1. Provider Interface

Create a common interface for all LLM providers:

```go
type LLMProvider interface {
    Chat(ctx context.Context, messages []Message, opts ChatOptions) (Message, error)
    StreamChat(ctx context.Context, messages []Message, opts ChatOptions, callback func(Message)) error
    SupportedModels() []string
    Name() string
}

type ChatOptions struct {
    Temperature float64
    MaxTokens   int
    Stream      bool
    // ... provider-specific options
}
```

### 2. Separate Message Model

Consider moving message types to a separate package:

```
pkg/
  llm/
    messages/      # Pure message types
      message.go
      content.go
      tools.go
    providers/     # Provider implementations
      ollama/
      openai/
      anthropic/
    integration/   # Graph integration
      nodes.go
```

### 3. Builder Pattern for Nodes

More flexible node creation:

```go
node, err := ollama.NewChatNode("ChatNode").
    WithModel("llama2").
    WithTemperature(0.7).
    WithTimeout(30 * time.Second).
    WithSystemPrompt("You are a helpful assistant").
    Build()
```

---

## 📚 Documentation Needs

- [ ] **README.md for pkg/llm**
  - Overview of message model
  - Supported providers
  - Basic usage examples
  - Provider comparison table

- [ ] **Provider-specific READMEs**
  - `pkg/llm/ollama/README.md`
  - `pkg/llm/openai/README.md`
  - `pkg/llm/aiw/README.md`

- [ ] **Examples**
  - Basic chat bot
  - Multi-turn conversation
  - Function calling
  - RAG pattern
  - Multi-provider fallback

- [ ] **API Documentation**
  - Godoc for all exported types
  - Usage examples in godoc
  - Best practices guide

---

## 🔍 Known Issues

1. **🔥 AIW Base URL Hack** (`aiw/aiw_chat.go:18`) - **CRITICAL**
   - Using `os.Setenv()` to set OpenAI base URL
   - Marked with TODO comment
   - **Breaks OpenAI client if multiple providers are used**
   - Affects global environment

2. **No Streaming in AIW**
   - OpenAI and Ollama have streaming
   - AIW uses synchronous call only

3. **Timestamp Always `time.Now()`**
   - Message timestamps are created on conversion
   - Loses original message timestamps from providers

4. **Partial Message Accumulation**
   - Complex logic in `ollama_chat.go:32-40`
   - Using `init` flag to track first message
   - Could be cleaner with proper state management

5. **No Rate Limiting**
   - No built-in rate limiting for any provider
   - Could hit API rate limits

---

## ✅ Completed / Working Well

- ✅ Basic chat functionality for Ollama, OpenAI, AIW
- ✅ Streaming support for Ollama and OpenAI
- ✅ Message role abstraction
- ✅ Basic message model with timestamps
- ✅ Integration with graph framework via builders
- ✅ Environment-based configuration

---

## 🎯 Recommended Priority Order

### Phase 1: Critical Fixes (This Week)
1. **Fix AIW base URL configuration** - Critical for multi-provider usage
2. Add context propagation to all providers
3. Implement proper streaming error handling

### Phase 2: Flexibility Improvements (Next 2 Weeks)
4. Improve node interfaces with options/builder patterns
5. Add configurable timeouts and provider options
6. Add basic unit tests

### Phase 3: Core Features (Next 1 Month)
7. Implement multimodal content support
8. Add tool/function calling support
9. Add helper methods to AgentModel
10. Add comprehensive tests

### Phase 4: Polish & Extend (Next 2-3 Months)
11. Add provider-specific enhancements
12. Add observability (logging, metrics)
13. Create detailed documentation
14. Add more provider implementations

### Phase 5: Advanced Features (Future)
15. Add additional providers (Anthropic, Gemini)
16. Implement conversation management
17. Add caching and retry logic
18. Build example applications

---

**Last Updated:** October 15, 2025  
**Maintainer:** ggraph team