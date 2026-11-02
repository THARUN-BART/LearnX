package aitutor

import (
	"fmt"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/ai"
)

const SystemPromptBase = `You are LearnX's Expert DSA Tutor and Systems Architect.
Your mission is to guide students to understand Data Structures & Algorithms deeply through Socratic questioning, structured conceptual breakdowns, and progressive scaffolding.

CRITICAL PEDAGOGICAL RULES:
1. NEVER immediately reveal full completed solution code unless explicitly requested at progressive hint Level 5.
2. In early stages (Hints 1-3), provide conceptual nudges, invariant hints, and edge cases. Ask a guiding question that helps the student think.
3. NEVER fabricate or disclose hidden test cases. You may only refer to public sample inputs or create simple conceptual examples.
4. Ground your explanations in DSA fundamentals (time/space complexity, loop invariants, state representations).
5. Always respond with clear, encouraging, precise technical language.
6. When structured JSON output is requested, output valid RFC 8259 JSON ONLY, without markdown code blocks, backticks, or prefix text.
`

func BuildHintMessages(level int, problemTitle, description, constraints string, samples string, studentCode, lang string) []ai.ChatMessage {
	var levelInstruction string
	switch level {
	case HintLevelConcept:
		levelInstruction = `LEVEL 1: CONCEPT HINT
- Explain the underlying DSA concept or strategy (e.g., Two Pointers, Monotonic Stack, Sliding Window, DP State).
- DO NOT provide any code or pseudocode.
- Provide a thought-provoking guiding question to help the student formulate the invariant.`
	case HintLevelStructural:
		levelInstruction = `LEVEL 2: STRUCTURAL / APPROACH HINT
- Outline the high-level steps or invariants needed (e.g., how pointers move, what map keys represent).
- Mention any state that must be maintained.
- DO NOT write the full code. Include a small pseudo-structure or guiding question.`
	case HintLevelEdgeCase:
		levelInstruction = `LEVEL 3: EDGE CASE HINT
- Identify critical edge cases or boundary conditions (empty inputs, duplicates, single elements, negative numbers, integer overflow).
- Ask how the current approach handles these edge cases.`
	case HintLevelBlueprint:
		levelInstruction = `LEVEL 4: ALGORITHMIC BLUEPRINT
- Provide structured pseudocode or step-by-step algorithmic blueprint.
- Show the loop structures and branch logic clearly without giving language-specific boilerplate.`
	case HintLevelCodeFix:
		levelInstruction = `LEVEL 5: CODE WALKTHROUGH & SOLUTION FIX
- In language: ` + lang + `.
- Point out the exact bug or missing logic in the student's attempt.
- Provide the idiomatic, fully commented solution and explain why it satisfies the time and space complexity.`
	default:
		levelInstruction = `Provide a helpful progressive hint for the student.`
	}

	system := SystemPromptBase + "\n\n" + levelInstruction + `
Respond in JSON format with the following keys:
{
  "title": "<Short descriptive title>",
  "content": "<Detailed hint text>",
  "guiding_question": "<A thought-provoking question for the student>"
}
`

	userContent := fmt.Sprintf(`PROBLEM: %s
DESCRIPTION: %s
CONSTRAINTS: %s
SAMPLE TEST CASES:
%s

STUDENT CURRENT CODE (%s):
%s`, problemTitle, description, constraints, samples, lang, studentCode)

	return []ai.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userContent},
	}
}

func BuildExplainMistakeMessages(problemTitle string, category, subtype, diagDesc, correction, errOutput, code, lang string) []ai.ChatMessage {
	system := SystemPromptBase + `
TASK: Explain the student's submission mistake in friendly, instructive pedagogical terms.
Analyze why this error happened and how they can master this concept.

Respond in JSON format with keys:
{
  "what_went_wrong": "<Clear explanation of the error>",
  "why_it_happened": "<Underlying cause or conceptual misconception>",
  "how_to_fix": "<Concrete suggestions and logic changes to resolve it>",
  "related_dsa_concepts": ["<Concept 1>", "<Concept 2>"]
}
`
	userContent := fmt.Sprintf(`PROBLEM: %s
MISTAKE CATEGORY: %s
SUBTYPE: %s
DIAGNOSIS SUMMARY: %s
SUGGESTED CORRECTION: %s
EXECUTION ERROR OUTPUT:
%s

STUDENT CODE (%s):
%s`, problemTitle, category, subtype, diagDesc, correction, errOutput, lang, code)

	return []ai.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userContent},
	}
}

func BuildWalkthroughMessages(problemTitle string, code, lang string) []ai.ChatMessage {
	system := SystemPromptBase + `
TASK: Conduct a comprehensive code walkthrough of the student's DSA solution.
Break down each major section of the code, noting strengths, potential pitfalls, and edge cases.

Respond in JSON format with keys:
{
  "overall_summary": "<Summary of the solution approach and architecture>",
  "sections": [
    {
      "title": "<Section Name, e.g. Initialization / Main Loop>",
      "line_range": "<e.g. Lines 5-12>",
      "observation": "<What this section does and how it behaves>",
      "evaluation": "<'positive', 'warning', or 'critical'>"
    }
  ],
  "key_strengths": ["<Strength 1>", "<Strength 2>"],
  "areas_to_improve": ["<Area 1>", "<Area 2>"]
}
`
	userContent := fmt.Sprintf(`PROBLEM: %s
LANGUAGE: %s
SOURCE CODE:
%s`, problemTitle, lang, code)

	return []ai.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userContent},
	}
}

func BuildDryRunMessages(problemTitle string, inputTestCase string, code, lang string) []ai.ChatMessage {
	system := SystemPromptBase + `
TASK: Generate a step-by-step dry-run simulation of the provided DSA code for the frontend visualizer.
Trace variable changes, data structure states, loop counters, and conditions.

Respond in JSON format with keys:
{
  "test_case_input": "<Input string>",
  "steps": [
    {
      "step_number": 1,
      "line_number": <int>,
      "action": "<Short description of action>",
      "variables": { "varName": "val" },
      "explanation": "<Why this step executed and what changed>"
    }
  ],
  "final_output": "<String representation of output>"
}
`
	userContent := fmt.Sprintf(`PROBLEM: %s
TEST CASE INPUT: %s
LANGUAGE: %s
SOURCE CODE:
%s`, problemTitle, inputTestCase, lang, code)

	return []ai.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userContent},
	}
}

func BuildComplexityMessages(problemTitle, constraints string, code, lang string) []ai.ChatMessage {
	system := SystemPromptBase + `
TASK: Analyze the Time and Space Complexity of the provided DSA code.
Compare it against optimal bounds based on typical competitive programming constraints.

Respond in JSON format with keys:
{
  "time_complexity": "<e.g. O(N) or O(N log N)>",
  "time_explanation": "<Detailed mathematical/loop explanation>",
  "space_complexity": "<e.g. O(1) or O(N)>",
  "space_explanation": "<Explanation of auxiliary memory usage>",
  "is_optimal": <true or false>,
  "optimisation_notes": "<Suggestions for reaching optimal complexity if not already optimal>"
}
`
	userContent := fmt.Sprintf(`PROBLEM: %s
CONSTRAINTS: %s
LANGUAGE: %s
SOURCE CODE:
%s`, problemTitle, constraints, lang, code)

	return []ai.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userContent},
	}
}

func BuildAskMessages(history []ai.ChatMessage, question string, problemContext string) []ai.ChatMessage {
	system := SystemPromptBase + `
You are in conversational tutoring mode. Answer the student's question clearly with Socratic guidance, analogies, and intuition.
If asked for code, guide them through reasoning rather than dumping complete answers immediately, unless they are truly stuck or explicitly asked for syntax clarification.
`
	if problemContext != "" {
		system += "\nCURRENT PROBLEM CONTEXT:\n" + problemContext
	}

	messages := []ai.ChatMessage{
		{Role: "system", Content: system},
	}

	// Append previous messages up to 10 for context window efficiency
	startIdx := 0
	if len(history) > 10 {
		startIdx = len(history) - 10
	}
	messages = append(messages, history[startIdx:]...)

	messages = append(messages, ai.ChatMessage{
		Role:    "user",
		Content: strings.TrimSpace(question),
	})

	return messages
}
