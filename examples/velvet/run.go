package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	a "github.com/morphy76/ggraph/pkg/agent"
	aiw "github.com/morphy76/ggraph/pkg/agent/aiw"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// TeacherNodeFn creates a conversational node for a high school teacher generating questions
var TeacherNodeFn o.ConversationNodeFn = func(chatService openai.ChatService, model string, opts ...option.RequestOption) g.NodeFn[a.Conversation] {
	return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		// System instruction for the teacher
		systemMsg := a.CreateMessage(a.System, "Sei un insegnante di scuola superiore. "+
			"Genera una domanda casuale su cultura generale, matematica, fisica, letteratura, scienze, storia, in generale su argomenti scolastici ma non di attualità. "+
			"Fai la domanda, con una brevissima introduzione o spiegazione. Parla solo in italiano.")

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMsg.Content),
			openai.UserMessage("Genera una domanda casuale."),
		}

		resp, err := chatService.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: messages,
		})
		if err != nil {
			return currentState, fmt.Errorf("failed to generate question: %w", err)
		}

		// Add the teacher's question to the conversation
		question := resp.Choices[0].Message.Content
		currentState.Messages = append(currentState.Messages,
			a.CreateMessage(a.Assistant, question))

		return currentState, nil
	}
}

// StudentNodeFn creates a conversational node for a student answering questions
var StudentNodeFn o.ConversationNodeFn = func(chatService openai.ChatService, model string, opts ...option.RequestOption) g.NodeFn[a.Conversation] {
	return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		// Get the last message (the question from the teacher)
		if len(currentState.Messages) == 0 {
			return currentState, fmt.Errorf("no question to answer")
		}

		lastMessage := currentState.Messages[len(currentState.Messages)-1]
		question := lastMessage.Content

		// System instruction for the student
		systemMsg := a.CreateMessage(a.System, "Sei uno studente di scuola superiore. "+
			"Rispondi alla domanda che ti viene posta nel modo più completo e preciso possibile. "+
			"Parla solo in italiano.")

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMsg.Content),
			openai.UserMessage(question),
		}

		resp, err := chatService.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: messages,
		})
		if err != nil {
			return currentState, fmt.Errorf("failed to generate answer: %w", err)
		}

		// Add the student's answer to the conversation
		answer := resp.Choices[0].Message.Content
		currentState.Messages = append(currentState.Messages,
			a.CreateMessage(a.User, answer))

		return currentState, nil
	}
}

// EvaluatorNodeFn creates a conversational node for an expert linguist evaluating answers
var EvaluatorNodeFn o.ConversationNodeFn = func(chatService openai.ChatService, model string, opts ...option.RequestOption) g.NodeFn[a.Conversation] {
	return func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		// Get the question and answer
		if len(currentState.Messages) < 2 {
			return currentState, fmt.Errorf("not enough messages to evaluate")
		}

		question := currentState.Messages[len(currentState.Messages)-2].Content
		answer := currentState.Messages[len(currentState.Messages)-1].Content

		// System instruction for the evaluator
		systemMsg := a.CreateMessage(a.System, "Sei un esperto linguista. "+
			"Valuta la risposta dello studente in termini di grammatica e correttezza lessicale. "+
			"Dai un punteggio da 0 a 10 per ogni categoria (10 è il massimo) e un breve commento. "+
			"Rispondi SOLO con un oggetto JSON nel seguente formato: "+
			`{"grammatica": {"punteggio": <numero>, "commento": "<testo>"}, "lessico": {"punteggio": <numero>, "commento": "<testo>"}}. `+
			"Non aggiungere altro testo prima o dopo il JSON. Parla solo in italiano nei commenti.")

		prompt := fmt.Sprintf("Domanda: %s\n\nRisposta: %s\n\nValuta la risposta.", question, answer)

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMsg.Content),
			openai.UserMessage(prompt),
		}

		resp, err := chatService.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: messages,
		})
		if err != nil {
			return currentState, fmt.Errorf("failed to generate evaluation: %w", err)
		}

		// Add the evaluator's assessment to the conversation
		evaluation := resp.Choices[0].Message.Content
		currentState.Messages = append(currentState.Messages,
			a.CreateMessage(a.Assistant, evaluation))

		return currentState, nil
	}
}

func main() {
	// Get AIW API key from environment
	pat := aiw.PATFromEnv()
	if pat == "" {
		log.Fatal("AIW_API_KEY environment variable not set")
	}

	// Create the three nodes with different Velvet models
	teacherNode, err := aiw.CreateConversationNode(
		"TeacherNode",
		pat,
		"velvet-2b",
		TeacherNodeFn,
	)
	if err != nil {
		log.Fatalf("Failed to create teacher node: %v", err)
	}

	studentNode, err := aiw.CreateConversationNode(
		"StudentNode",
		pat,
		"velvet-25b-07-15771",
		StudentNodeFn,
	)
	if err != nil {
		log.Fatalf("Failed to create student node: %v", err)
	}

	evaluatorNode, err := aiw.CreateConversationNode(
		"EvaluatorNode",
		pat,
		"velvet-14b",
		EvaluatorNodeFn,
	)
	if err != nil {
		log.Fatalf("Failed to create evaluator node: %v", err)
	}

	// Create edges connecting the nodes
	startEdge := b.CreateStartEdge(teacherNode)
	teacherToStudentEdge := b.CreateEdge(teacherNode, studentNode)
	studentToEvaluatorEdge := b.CreateEdge(studentNode, evaluatorNode)
	endEdge := b.CreateEndEdge(evaluatorNode)

	// Initialize the conversation state
	initialState := a.CreateConversation()
	stateMonitorCh := make(chan g.StateMonitorEntry[a.Conversation], 10)

	// Create the runtime graph
	graph, err := b.CreateRuntimeWithInitialState(startEdge, stateMonitorCh, initialState)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()

	// Add all edges to the graph
	graph.AddEdge(teacherToStudentEdge, studentToEvaluatorEdge, endEdge)

	// Validate the graph
	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	// Run the graph with empty user input (no user interaction needed)
	userInput := a.CreateConversation()
	threadID := graph.Invoke(userInput)

	// Monitor and display the conversation flow
	fmt.Println("=== Velvet Educational Example ===")
	fmt.Println("Processo: Insegnante → Studente → Valutatore")
	fmt.Println()

	for {
		entry := <-stateMonitorCh

		if entry.Running {
			fmt.Printf("▶ Esecuzione nodo: %s\n", entry.Node)
		} else {
			if entry.Error != nil {
				fmt.Printf("✗ Errore nel nodo %s: %v\n\n", entry.Node, entry.Error)
			} else {
				fmt.Printf("✓ Completato nodo: %s\n", entry.Node)

				// Display the latest message
				if len(entry.NewState.Messages) > 0 {
					lastMsg := entry.NewState.Messages[len(entry.NewState.Messages)-1]

					switch entry.Node {
					case "TeacherNode":
						fmt.Printf("\n📚 DOMANDA DELL'INSEGNANTE:\n%s\n\n", lastMsg.Content)
					case "StudentNode":
						fmt.Printf("\n🎓 RISPOSTA DELLO STUDENTE:\n%s\n\n", lastMsg.Content)
					case "EvaluatorNode":
						fmt.Printf("\n📊 VALUTAZIONE LINGUISTICA:\n")
						// Try to parse and pretty-print the JSON
						var evaluation map[string]interface{}
						if err := json.Unmarshal([]byte(lastMsg.Content), &evaluation); err == nil {
							prettyJSON, _ := json.MarshalIndent(evaluation, "", "  ")
							fmt.Printf("%s\n\n", prettyJSON)
						} else {
							fmt.Printf("%s\n\n", lastMsg.Content)
						}
					}
				}
			}
		}

		if !entry.Running {
			// Pretty print the final graph status
			state := graph.CurrentState(threadID)
			statusJSON, _ := json.MarshalIndent(state, "", "  ")
			fmt.Printf("\n📈 STATO FINALE DEL GRAFO:\n%s\n\n", statusJSON)
			break
		}
	}

	fmt.Println("=== Fine del processo ===")
}
