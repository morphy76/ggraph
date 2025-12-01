package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	a "github.com/morphy76/ggraph/pkg/agent"
	aiw "github.com/morphy76/ggraph/pkg/agent/aiw"
	o "github.com/morphy76/ggraph/pkg/agent/openai"
	b "github.com/morphy76/ggraph/pkg/builders"
	g "github.com/morphy76/ggraph/pkg/graph"
)

// Evaluation represents a parsed evaluation result
type Evaluation struct {
	Grammatica struct {
		Punteggio int    `json:"punteggio"`
		Commento  string `json:"commento"`
	} `json:"grammatica"`
	Lessico struct {
		Punteggio int    `json:"punteggio"`
		Commento  string `json:"commento"`
	} `json:"lessico"`
	Contenuto struct {
		Punteggio int    `json:"punteggio"`
		Commento  string `json:"commento"`
	} `json:"contenuto"`
}

// ThreadResult holds the results of a single thread execution
type ThreadResult struct {
	ThreadID   string
	Success    bool
	Evaluation *Evaluation
	Error      error
	StartTime  time.Time
	EndTime    time.Time
}

// ThreadMonitor tracks state for a single thread
type ThreadMonitor struct {
	threadNum       int
	threadID        string
	finalEvaluation *Evaluation
	lastError       error
	messages        []string
	completed       bool
	lastNode        string
	seenEvaluator   bool
}

// runThread executes a single thread and monitors its progress
func runThread(
	graph g.Runtime[a.Conversation],
	threadNum int,
	resultChan chan<- ThreadResult,
	progressMutex *sync.Mutex,
	threadProgress map[int]string,
	monitors map[string]*ThreadMonitor,
	verbose bool,
) {
	startTime := time.Now()
	result := ThreadResult{
		ThreadID:  fmt.Sprintf("Thread-%d", threadNum),
		StartTime: startTime,
	}

	// Update progress
	updateProgress := func(status string) {
		progressMutex.Lock()
		threadProgress[threadNum] = status
		progressMutex.Unlock()
	}

	updateProgress("Avvio")

	// Pre-assign a threadID to avoid race condition
	threadID := fmt.Sprintf("thread-%d", threadNum)

	// Create and register monitor BEFORE invoking the graph
	monitor := &ThreadMonitor{
		threadNum: threadNum,
		threadID:  threadID,
		messages:  make([]string, 0),
	}

	progressMutex.Lock()
	monitors[threadID] = monitor
	progressMutex.Unlock()

	if verbose {
		fmt.Printf("\n[Thread %d] Registrato con ID: %s\n", threadNum, threadID)
	}

	// Run the graph with the pre-assigned threadID
	userInput := a.CreateConversation(a.CreateMessage(a.User, "Genera una domanda casuale su cultura generale, matematica, fisica, letteratura, scienze, storia, in generale su argomenti scolastici ma non di attualità."))
	actualThreadID := graph.Invoke(userInput, g.InvokeConfigThreadID(threadID))

	fmt.Printf("\n✅ [Thread %d] Invoke completato: actualThreadID=%s, expected=%s\n",
		threadNum, actualThreadID, threadID)

	if actualThreadID != threadID {
		fmt.Printf("\n⚠️  [Thread %d] ThreadID DIVERSO da atteso!\n", threadNum)
	}

	// Wait for this specific thread to complete
	timeout := time.After(120 * time.Second)
	checkInterval := time.NewTicker(100 * time.Millisecond)
	defer checkInterval.Stop()

	checkCount := 0
	for {
		select {
		case <-timeout:
			updateProgress("Timeout")

			// Final debug before timeout
			progressMutex.Lock()
			fmt.Printf("\n❌ [Thread %d] TIMEOUT! Status finale:\n", threadNum)
			fmt.Printf("   - completed: %v\n", monitor.completed)
			fmt.Printf("   - finalEvaluation: %v\n", monitor.finalEvaluation != nil)
			fmt.Printf("   - lastError: %v\n", monitor.lastError)
			fmt.Printf("   - messages count: %d\n", len(monitor.messages))
			progressMutex.Unlock()

			result.EndTime = time.Now()
			result.Success = false
			result.Error = fmt.Errorf("timeout dopo 120 secondi")
			resultChan <- result
			return

		case <-checkInterval.C:
			checkCount++
			progressMutex.Lock()
			completed := monitor.completed
			finalEval := monitor.finalEvaluation
			lastErr := monitor.lastError

			// Debug every 10 seconds
			if verbose && checkCount%100 == 0 {
				fmt.Printf("\n🔍 [Thread %d] Check #%d: completed=%v, hasEval=%v, hasError=%v\n",
					threadNum, checkCount, completed, finalEval != nil, lastErr != nil)
			}
			progressMutex.Unlock()

			if completed {
				result.EndTime = time.Now()
				result.Success = lastErr == nil && finalEval != nil
				result.Evaluation = finalEval
				result.Error = lastErr

				if verbose && result.Success {
					updateProgress("✓ Completato")
				} else if verbose && !result.Success {
					updateProgress(fmt.Sprintf("✗ Errore: %v", lastErr))
				} else {
					updateProgress("Completato")
				}

				resultChan <- result
				return
			}
		}
	}
}

// calculateAverages computes the average scores from all successful evaluations
func calculateAverages(results []ThreadResult) (avgGrammatica, avgLessico, avgContenuto float64, successCount int) {
	var totalGrammatica, totalLessico, totalContenuto int

	for _, result := range results {
		if result.Success && result.Evaluation != nil {
			totalGrammatica += result.Evaluation.Grammatica.Punteggio
			totalLessico += result.Evaluation.Lessico.Punteggio
			totalContenuto += result.Evaluation.Contenuto.Punteggio
			successCount++
		}
	}

	if successCount > 0 {
		avgGrammatica = float64(totalGrammatica) / float64(successCount)
		avgLessico = float64(totalLessico) / float64(successCount)
		avgContenuto = float64(totalContenuto) / float64(successCount)
	}

	return
}

// monitorStateChannel processes all state updates from the graph
func monitorStateChannel(
	stateMonitorCh <-chan g.StateMonitorEntry[a.Conversation],
	monitors map[string]*ThreadMonitor,
	progressMutex *sync.Mutex,
	threadProgress map[int]string,
	verbose bool,
	done chan bool,
) {
	eventCount := 0
	for entry := range stateMonitorCh {
		eventCount++

		if verbose {
			fmt.Printf("\n🔔 Evento #%d ricevuto: ThreadID=%s, Node=%s, Running=%v, Error=%v\n",
				eventCount, entry.ThreadID, entry.Node, entry.Running, entry.Error)
		}

		progressMutex.Lock()

		monitor, exists := monitors[entry.ThreadID]
		if !exists {
			fmt.Printf("\n⚠️  Evento #%d: ThreadID sconosciuto: %s (Node: %s, Running: %v)\n",
				eventCount, entry.ThreadID, entry.Node, entry.Running)
			fmt.Printf("   Monitor registrati: %d\n", len(monitors))
			if len(monitors) > 0 {
				fmt.Printf("   ThreadID registrati: ")
				for tid := range monitors {
					fmt.Printf("%s ", tid)
				}
				fmt.Println()
			}
			progressMutex.Unlock()
			continue
		}

		// Track which node we're on
		monitor.lastNode = entry.Node
		if entry.Node == "EvaluatorNode" {
			monitor.seenEvaluator = true
		}

		if entry.Running {
			// Thread is still running - update progress
			threadProgress[monitor.threadNum] = fmt.Sprintf("Esecuzione: %s", entry.Node)

			if verbose {
				fmt.Printf("\n[Thread %d] ▶ %s (thread attivo)\n", monitor.threadNum, entry.Node)
			}
		} else {
			// Thread has completed (Running=false)
			if entry.Error != nil {
				monitor.lastError = entry.Error
				threadProgress[monitor.threadNum] = fmt.Sprintf("Errore: %v", entry.Error)

				if verbose {
					fmt.Printf("\n[Thread %d] ✗ Thread completato con errore: %v\n", monitor.threadNum, entry.Error)
				}
			} else {
				threadProgress[monitor.threadNum] = "Thread completato"

				if verbose {
					fmt.Printf("\n[Thread %d] ✅ Thread completato con successo\n", monitor.threadNum)
				}

				// Extract final evaluation from messages (should have 3: question, answer, evaluation)
				if len(entry.NewState.Messages) >= 3 {
					evalMsg := entry.NewState.Messages[2] // Third message should be evaluation
					var eval Evaluation
					if err := json.Unmarshal([]byte(evalMsg.Content), &eval); err == nil {
						monitor.finalEvaluation = &eval
						if verbose {
							fmt.Printf("\n[Thread %d] 📊 Valutazione finale estratta:\n", monitor.threadNum)
							fmt.Printf("   Grammatica: %d/10 - %s\n", eval.Grammatica.Punteggio, eval.Grammatica.Commento)
							fmt.Printf("   Lessico: %d/10 - %s\n", eval.Lessico.Punteggio, eval.Lessico.Commento)
							fmt.Printf("   Contenuto: %d/10 - %s\n", eval.Contenuto.Punteggio, eval.Contenuto.Commento)
						}
					} else if verbose {
						fmt.Printf("\n[Thread %d] ⚠️  Errore parsing valutazione finale: %v\n", monitor.threadNum, err)
						fmt.Printf("   Contenuto messaggio #3: %.200s\n", evalMsg.Content)
					}

					// Display all messages
					if verbose {
						fmt.Printf("\n[Thread %d] 📚 DOMANDA:\n%s\n", monitor.threadNum, entry.NewState.Messages[0].Content)
						fmt.Printf("\n[Thread %d] 🎓 RISPOSTA:\n%s\n", monitor.threadNum, entry.NewState.Messages[1].Content)
						fmt.Printf("\n[Thread %d] 📊 VALUTAZIONE:\n%s\n", monitor.threadNum, entry.NewState.Messages[2].Content)
					}
				} else if verbose {
					fmt.Printf("\n[Thread %d] ⚠️  Messaggi insufficienti: %d (attesi 3)\n",
						monitor.threadNum, len(entry.NewState.Messages))
				}
			}

			// Mark thread as completed
			monitor.completed = true
		}

		progressMutex.Unlock()
	}

	done <- true
}

// printProgress displays the current progress of all threads
func printProgress(threadProgress map[int]string, totalThreads int) {
	fmt.Print("\033[2J\033[H") // Clear screen and move to top
	fmt.Println("=== Velvet Educational Example - Esecuzione Concorrente ===")
	fmt.Println()
	fmt.Printf("Thread attivi: %d\n", totalThreads)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i := 1; i <= totalThreads; i++ {
		status := threadProgress[i]
		if status == "" {
			status = "In attesa"
		}
		fmt.Printf("Thread %2d: %s\n", i, status)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func main() {
	// Get user input for number of threads and time period
	var numThreads int
	var timePeriod int

	fmt.Println("=== Velvet Educational Example - Configurazione ===")
	fmt.Print("Inserisci il numero di thread concorrenti (max 20): ")
	if _, err := fmt.Scanf("%d", &numThreads); err != nil || numThreads < 1 || numThreads > 20 {
		log.Fatal("Numero di thread non valido (deve essere tra 1 e 20)")
	}

	fmt.Print("Inserisci il periodo temporale in secondi (max 120): ")
	if _, err := fmt.Scanf("%d", &timePeriod); err != nil || timePeriod < 1 || timePeriod > 120 {
		log.Fatal("Periodo temporale non valido (deve essere tra 1 e 120)")
	}

	// Get AIW API key from environment
	pat := aiw.PATFromEnv()
	if pat == "" {
		log.Fatal("AIW_API_KEY environment variable not set; visit https://portal.aiwave.ai to get your API key.")
	}

	aiwClient := aiw.NewAIWClient(pat)

	// Create the three nodes with different Velvet models
	teacherNode, err := o.CreateConversationNode(
		"TeacherNode",
		"velvet-2b-1.5-03-23918",
		aiwClient,
		a.WithMessages(
			a.CreateMessage(a.System, "Sei un insegnante di scuola superiore che sta interrogando uno studente poco preparato. "+
				"Non dare suggerimenti o risposte alla domanda posta. Fai la domanda, con una brevissima introduzione o spiegazione. "+
				"Parla solo in italiano."),
		),
		a.WithTemperature(1.0),
	)
	if err != nil {
		log.Fatalf("Failed to create teacher node: %v", err)
	}
	askNode, err := b.NewNode("Teach2Student", func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		return a.CreateConversation(
			currentState.Messages[len(currentState.Messages)-1],
		), nil
	})

	studentNode, err := o.CreateConversationNode(
		"StudentNode",
		"velvet-25b-07-15771",
		aiwClient,
		a.WithMessages(
			a.CreateMessage(a.System, "Sei uno studente di scuola superiore. "+
				"Rispondi alla domanda che ti viene posta nel modo più completo e preciso possibile. "+
				"Parla solo in italiano."),
		),
		a.WithTemperature(0.5),
	)
	if err != nil {
		log.Fatalf("Failed to create student node: %v", err)
	}
	answerNode, err := b.NewNode("Student2Eval", func(userInput, currentState a.Conversation, notify g.NotifyPartialFn[a.Conversation]) (a.Conversation, error) {
		mexCount := len(currentState.Messages)
		return a.CreateConversation(
			a.CreateMessage(a.User, fmt.Sprintf("Question:\n%s\n\nAnswer:\n%s", currentState.Messages[mexCount-2].Content, currentState.Messages[mexCount-1].Content)),
		), nil
	})

	evaluatorNode, err := o.CreateConversationNode(
		"EvaluatorNode",
		"velvet-14b",
		aiwClient,
		a.WithMessages(
			a.CreateMessage(a.System, "Sei un esperto linguista e valutatore di contenuti. "+
				"Valuta la risposta dello studente in termini di grammatica, correttezza lessicale e correttezza del contenuto rispetto alla domanda posta. "+
				"Dai un punteggio da 0 a 10 per ogni categoria (10 è il massimo) e un breve commento. "+
				"Rispondi SOLO con un oggetto JSON nel seguente formato: "+
				`{"grammatica": {"punteggio": <numero>, "commento": "<testo>"}, "lessico": {"punteggio": <numero>, "commento": "<testo>"}, "contenuto": {"punteggio": <numero>, "commento": "<testo>"}}. `+
				"Non aggiungere altro testo prima o dopo il JSON. Parla solo in italiano nei commenti."),
		),
		a.WithTemperature(0.0),
	)
	if err != nil {
		log.Fatalf("Failed to create evaluator node: %v", err)
	}

	// Create edges connecting the nodes
	startEdge := b.CreateStartEdge(teacherNode)
	teacher2Ask := b.CreateEdge(teacherNode, askNode)
	ask2StudentEdge := b.CreateEdge(askNode, studentNode)
	student2Answer := b.CreateEdge(studentNode, answerNode)
	answer2EvaluatorEdge := b.CreateEdge(answerNode, evaluatorNode)
	endEdge := b.CreateEndEdge(evaluatorNode)

	// Initialize the conversation state
	initialState := a.CreateConversation()
	stateMonitorCh := make(chan g.StateMonitorEntry[a.Conversation], numThreads*10)

	// Create the runtime graph
	graph, err := b.CreateRuntime(startEdge, stateMonitorCh, g.WithInitialState(initialState))
	if err != nil {
		log.Fatalf("Runtime creation failed: %v", err)
	}
	defer graph.Shutdown()

	// Add all edges to the graph
	graph.AddEdge(teacher2Ask, ask2StudentEdge, student2Answer, answer2EvaluatorEdge, endEdge)

	// Validate the graph
	err = graph.Validate()
	if err != nil {
		log.Fatalf("Graph validation failed: %v", err)
	}

	// Setup concurrent execution
	resultChan := make(chan ThreadResult, numThreads)
	var progressMutex sync.Mutex
	threadProgress := make(map[int]string)
	monitors := make(map[string]*ThreadMonitor)
	var wg sync.WaitGroup

	// Ask for verbose mode
	var verboseInput string
	fmt.Print("Modalità verbose (mostra domande/risposte/valutazioni)? (s/n): ")
	fmt.Scanf("%s", &verboseInput)
	verbose := verboseInput == "s" || verboseInput == "S"

	// Calculate delay between thread starts
	delayBetweenThreads := time.Duration(timePeriod) * time.Second / time.Duration(numThreads)

	fmt.Println()
	fmt.Printf("Avvio di %d thread distribuiti su %d secondi (intervallo: %.2f secondi)\n",
		numThreads, timePeriod, delayBetweenThreads.Seconds())
	if !verbose {
		fmt.Println("Premi Invio per iniziare...")
		fmt.Scanln()
	} else {
		var dummy string
		fmt.Println("Premi Invio per iniziare...")
		fmt.Scanln(&dummy)
		fmt.Println()
		fmt.Println("=== INIZIO ESECUZIONE ===")
		fmt.Println()
	}

	startTime := time.Now()

	// Start global state monitor
	monitorDone := make(chan bool)
	go monitorStateChannel(stateMonitorCh, monitors, &progressMutex, threadProgress, verbose, monitorDone)

	fmt.Println()
	fmt.Printf("✅ Monitor globale avviato (canale buffer: %d)\n", cap(stateMonitorCh))
	fmt.Println()

	// Start progress monitoring (only in non-verbose mode)
	var progressTicker *time.Ticker
	if !verbose {
		progressTicker = time.NewTicker(500 * time.Millisecond)
		defer progressTicker.Stop()

		go func() {
			for range progressTicker.C {
				progressMutex.Lock()
				printProgress(threadProgress, numThreads)
				progressMutex.Unlock()
			}
		}()
	}

	// Launch threads with time distribution
	for i := 1; i <= numThreads; i++ {
		wg.Add(1)
		go func(threadNum int) {
			defer wg.Done()
			runThread(graph, threadNum, resultChan, &progressMutex, threadProgress, monitors, verbose)
		}(i)

		// Wait before starting the next thread (except for the last one)
		if i < numThreads {
			time.Sleep(delayBetweenThreads)
		}
	}

	// Wait for all threads to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var results []ThreadResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Stop progress ticker if it was started
	if progressTicker != nil {
		progressTicker.Stop()
	}

	totalDuration := time.Since(startTime)

	// Wait a bit for final state updates
	time.Sleep(500 * time.Millisecond)

	// Display final results
	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Println("=== Velvet Educational Example - Risultati Finali ===")
	fmt.Println()
	fmt.Printf("Tempo totale di esecuzione: %.2f secondi\n", totalDuration.Seconds())
	fmt.Printf("Thread completati: %d/%d\n", len(results), numThreads)
	fmt.Println()

	// Calculate averages
	avgGrammatica, avgLessico, avgContenuto, successCount := calculateAverages(results)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 MEDIA DELLE VALUTAZIONI")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Valutazioni riuscite: %d/%d (%.1f%%)\n",
		successCount, numThreads, float64(successCount)/float64(numThreads)*100)
	fmt.Println()

	if successCount > 0 {
		fmt.Printf("Grammatica: %.2f/10\n", avgGrammatica)
		fmt.Printf("Lessico:    %.2f/10\n", avgLessico)
		fmt.Printf("Contenuto:  %.2f/10\n", avgContenuto)
		fmt.Printf("Media tot.: %.2f/10\n", (avgGrammatica+avgLessico+avgContenuto)/3)
	} else {
		fmt.Println("Nessuna valutazione completata con successo")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 DETTAGLIO PER THREAD")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, result := range results {
		fmt.Printf("\nThread %d:\n", i+1)
		fmt.Printf("  Durata: %.2f secondi\n", result.EndTime.Sub(result.StartTime).Seconds())

		if result.Success && result.Evaluation != nil {
			fmt.Printf("  ✓ Successo\n")
			fmt.Printf("  Grammatica: %d/10 - %s\n",
				result.Evaluation.Grammatica.Punteggio,
				result.Evaluation.Grammatica.Commento)
			fmt.Printf("  Lessico:    %d/10 - %s\n",
				result.Evaluation.Lessico.Punteggio,
				result.Evaluation.Lessico.Commento)
			fmt.Printf("  Contenuto:  %d/10 - %s\n",
				result.Evaluation.Contenuto.Punteggio,
				result.Evaluation.Contenuto.Commento)
		} else {
			fmt.Printf("  ✗ Fallito")
			if result.Error != nil {
				fmt.Printf(": %v", result.Error)
			}
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println("=== Fine del processo ===")
}
