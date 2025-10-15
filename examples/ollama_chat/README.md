# Chat Loop Example - Interactive Cooking Assistant

This example demonstrates a **continuous chat loop** with Ollama LLM, maintaining conversation memory across iterations.

## Prompt

### Initial

Create an example using the other examples as inspiration which demonstrates a loop graph where:

- the graph loops until the end user types "exit"
- for each non exit message, chat with ollama using CreateOLLamaChatNodeFromEnvironment
- set envs of ollama host to local host
- model is Almawave/Velvet:2B
- preserve chat memory in the graph status
- initialize the graph using system message which states the chat is about cooking

### Follow-up

After the user input, it hangs, on ctrl+c:

go run ./examples/chat_loop/run.go
🍳 Cooking Chat Assistant
💡 Ask me anything about cooking! Type 'exit' to quit.

💬 You: qual è la ricetta del giorno?
fatal error: all goroutines are asleep - deadlock!

### Fix

I cannot see assistant answers:

go run ./examples/chat_loop/run.go
🍳 Cooking Chat Assistant
💡 Ask me anything about cooking! Type 'exit' to quit.

💬 You: qual è la ricetta del giorno?
🤖 Assistant: 

💬 You: come puoi aiutarmi?
🤖 Assistant: 

💬 You: exit

👋 Goodbye! Happy cooking!

## Architecture

Simple graph: **ChatNode → End**, invoked repeatedly in a loop.

The loop is external to the graph - the main function reads input, adds it to state, invokes the graph, waits for completion, displays the response, and repeats.

## Flow

1. Initialize graph with system message (cooking assistant persona)
2. **Main loop**:
   - Read user input
   - If "exit", quit
   - Add user message to state
   - Invoke graph with current state
   - ChatNode processes and generates response
   - Display assistant response
   - Repeat

## Key Features

- **Stateful chat memory**: All messages preserved across graph invocations
- **System message initialization**: Sets AI persona as cooking assistant  
- **Continuous loop**: External loop invokes graph until user types "exit"
- **Ollama integration**: Uses `CreateOLLamaChatNodeFromEnvironment`
- **State merger**: `MergeChatModels` preserves conversation history

## Prerequisites

1. Install and run Ollama locally:
   ```bash
   ollama serve
   ```

2. Pull the model:
   ```bash
   ollama pull Almawave/Velevet-2B
   ```

## Environment

The example sets:
```bash
OLLAMA_HOST=http://localhost:11434
```

## Run

```bash
go run ./examples/chat_loop/run.go
```

## Example Session

```
🍳 Cooking Chat Assistant
💡 Ask me anything about cooking! Type 'exit' to quit.

💬 You: How do I make pasta?
🤖 Assistant: Boil water with salt, add pasta, cook for 8-10 minutes...

💬 You: What about the sauce?
🤖 Assistant: For a simple tomato sauce, sauté garlic in olive oil...

💬 You: exit
👋 Goodbye! Happy cooking!
```
