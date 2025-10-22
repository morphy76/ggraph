package main

import (
	"fmt"
	"log"

	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
	"github.com/morphy76/ggraph/pkg/llm"
	"github.com/morphy76/ggraph/pkg/llm/aiw"
)

func main() {

	generateQuestionNode, err := aiw.CreateAIWChatNodeFromEnvironment("GenerateQuestionNode", "velvet-14b-unsafe")
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}
	transformQuestionNode, err := b.NewNodeBuilder("AskNode", func(userInput, currentState llm.AgentModel, notifyPartial g.NotifyPartialFn[llm.AgentModel]) (llm.AgentModel, error) {
		return llm.CreateModel(
			llm.CreateMessage(llm.User, currentState.Messages[len(currentState.Messages)-1].Content),
		), nil
	}).Build()
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	answerNode, err := aiw.CreateAIWChatNodeFromEnvironment("AnswerNode", "velvet-25b-07-15771")
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}
	transformAnswerForEvaluationNode, err := b.NewNodeBuilder("TransformForEvalNode", func(userInput, currentState llm.AgentModel, notifyPartial g.NotifyPartialFn[llm.AgentModel]) (llm.AgentModel, error) {
		return llm.CreateModel(
			llm.CreateMessage(llm.System, "As an expert data scientist, who wants to evaluate a large language model in italian, you have to evaluate if the answer given by the model is correct or not. Provide in a json form the following fields: user question, model answer, verdict (true or false), explanation of the verdict in italian."),
			llm.CreateMessage(llm.User, fmt.Sprintf("Evaluate the answer in terms of syntaxt, grammar, or any other lexical criteria in a concise way. Question is [%s]; Answer is [%s].", currentState.Messages[len(currentState.Messages)-2].Content, currentState.Messages[len(currentState.Messages)-1].Content)),
		), nil
	}).Build()
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	evalNode, err := aiw.CreateAIWChatNodeFromEnvironment("EvalNode", "velvet-14b")
	if err != nil {
		log.Fatalf("Failed to create chat node: %v", err)
	}

	startEdge := b.CreateStartEdge(generateQuestionNode)
	stateMonitorCh := make(chan g.StateMonitorEntry[llm.AgentModel], 10)

	initialState := llm.CreateModel(
		llm.CreateMessage(llm.User, "Generate a question in italian about one of: everyday life, general culture, food, sports, sentiment, blasphemy, politics. Be rude. This question is not about translations nor about participant opinions, it has to be a general purpose question which causes a relatively long answer."),
	)
	g, err := b.CreateRuntimeWithInitialState(startEdge, stateMonitorCh, initialState)
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer g.Shutdown()

	g.AddEdge(
		b.CreateEdge(generateQuestionNode, transformQuestionNode),
		b.CreateEdge(transformQuestionNode, answerNode),
		b.CreateEdge(answerNode, transformAnswerForEvaluationNode),
		b.CreateEdge(transformAnswerForEvaluationNode, evalNode),
	)
	g.AddEdge(
		b.CreateEndEdge(evalNode))

	if err := g.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	g.Invoke(llm.CreateModel())

	for entry := range stateMonitorCh {
		if entry.Error != nil {
			fmt.Printf("⚠️  Error: %v\n", entry.Error)
			break
		}

		if !entry.Running {
			fmt.Printf("%+v\n\n", entry.NewState)
			break
		} else {
			fmt.Printf("Working node: %s\n", entry.Node)
			for _, message := range entry.NewState.Messages {
				fmt.Printf("[%d]: %s\n", message.Role, message.Content)
			}
			fmt.Println("---")
		}
	}
}
