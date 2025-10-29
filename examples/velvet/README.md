# Velvet Example - Educational Q&A Evaluation

This example demonstrates a graph using AIW (Almawave) conversational nodes with different Velvet models to create an automated educational question-answer-evaluation flow.

## Overview

The graph consists of three nodes that work together:

1. **TeacherNode** (velvet-2b): Acts as a high school teacher generating random questions about general culture or high school topics
2. **StudentNode** (velvet-25b): Acts as a student answering the teacher's question
3. **EvaluatorNode** (velvet-14b): Acts as an expert linguist evaluating the student's answer

All nodes communicate in Italian 🇮🇹

## Flow

```
TeacherNode → StudentNode → EvaluatorNode
(generates    (answers       (evaluates
 question)     question)      answer)
```

## Evaluation Output

The evaluator provides a JSON-formatted assessment with:
- **Grammar score** (0-10): Grammatical correctness evaluation with a brief comment
- **Lexical score** (0-10): Lexical correctness evaluation with a brief comment

Example output:
```json
{
  "grammatica": {
    "punteggio": 9,
    "commento": "Frase ben strutturata con uso corretto dei tempi verbali"
  },
  "lessico": {
    "punteggio": 8,
    "commento": "Vocabolario appropriato con buona varietà lessicale"
  }
}
```

## Requirements

- AIW (Almawave) API access
- `AIW_API_KEY` environment variable set with your API key

## Running the Example

```bash
export AIW_API_KEY="your-api-key-here"
go run run.go
```

## Features

- **No user input required**: The graph runs autonomously
- **Multi-model architecture**: Uses three different Velvet model sizes optimized for different tasks
- **Structured output**: Evaluation results are formatted as JSON for easy parsing
- **Italian language**: All interactions are in Italian
- **Real-time monitoring**: Shows the progress of each node as it executes

## Notes

- The velvet-2b model is used for simple question generation
- The velvet-25b model provides comprehensive answers
- The velvet-14b model offers detailed linguistic analysis
- All models are accessed through the AIW Platform API
