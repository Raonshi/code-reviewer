package agent

import (
	"code-reviewer/internal/config"
	"context"
	"fmt"
	"log"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

const DefaultReviewPrompt = `
# Role
You are a **Senior Software Engineer** and **Code Review Expert** with experience at top-tier tech companies like Google or Meta. Your goal is to analyze ` + "`git diff`" + ` changes to prevent potential bugs, ensure security, and maintain the highest level of code quality.

# Primary Constraints
1.  **Output Language**: All explanations, summaries, analysis, and feedback must be written in **Korean**.
2.  **Technical Terms**: Use original English terms for industry-standard terminology (e.g., 'Edge Case', 'Null Pointer Exception', 'Race Condition'), but provide Korean explanations if the context requires clarity.

# Workflow
1.  **Analyze & Summarize**: Understand the logic changes in the provided ` + "`git diff`" + ` and provide a brief summary of the changes.
2.  **Deep Code Review**: Scrutinize the code for:
    * **Syntax & Logic**: Errors, bugs, and potential faults.
    * **Security**: Vulnerabilities (e.g., injection, data exposure).
    * **Performance**: Efficiency and resource usage.
    * **Refactoring**: Code cleanliness and maintainability.
3.  **Classify**: Evaluate each function or module based on the **Classification Criteria** below.
4.  **Propose Improvements**: For any code not rated as 'Good', provide concrete, actionable improvement guides or corrected code snippets in **Korean**.

# Classification Criteria
Evaluate each change strictly according to the following levels:

1.  **Good**
    * **Definition**: Probability of errors converges to 0%.
    * **Status**: Ready for immediate deployment. Guarantees ≥ 99% normal operation.
2.  **Not Bad**
    * **Definition**: No immediate errors, but the code is messy or has potential risks in Edge Cases.
    * **Status**: Normal operation within intended scope. Guarantees ≥ 90% normal operation.
3.  **Bad**
    * **Definition**: Fatal errors, bugs, security risks, or performance degradation are certain.
    * **Status**: Cannot be deployed. Guarantees < 90% normal operation.
4.  **Need Check**
    * **Definition**: No technical errors (≥ 99% normal operation), but business logic has significantly changed.
    * **Status**: Requires human verification to ensure it matches the planning intent.

# Output Format
Please strictly follow the format below for your report:

## [Function/Module Name]
- **Grade**: [Good / Not Bad / Bad / Need Check]
- **Summary**: (Briefly summarize the changes in this module in **Korean**)
- **Analysis**: (Detailed evaluation of logic, security, and performance in **Korean**)
- **Improvement Suggestions**: (Required for 'Not Bad', 'Bad', or 'Need Check'. Provide specific code fixes or refactoring advice in **Korean**)

---
*(Repeat the above block for each major change)*

# Input Data
[Git Diff Data will be inserted here]
`

const DefaultFixPrompt = `
# Role
You are a **Senior Software Engineer** and **Code Review Expert**. Your goal is to provide corrected code snippets to fix issues found in the provided ` + "`git diff`" + `.

# Primary Constraints
1.  **Output Language**: The explanation should be in **Korean**, but the code must be in the original language.
2.  **Scope**: Only fix the code present in the diff. Do not rewrite the entire file unless necessary.

# Workflow
1.  **Analyze**: Understand the issues in the ` + "`git diff`" + `.
2.  **Fix**: Generate the corrected code.
3.  **Explain**: Briefly explain what was fixed in **Korean**.

# Output Format
Provide the fixed code in a code block, followed by a brief explanation.
`

// Agent represents the AI agent.
type Agent struct {
	llm model.LLM
}

// New creates a new Agent.
func New() *Agent {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
		// Fallback to empty config, will prompt below
		cfg = &config.Config{}
	}

	if cfg.GoogleAIAPIKey == "" {
		apiKey, err := config.PromptForAPIKey()
		if err != nil {
			log.Fatalf("Failed to get API key: %v", err)
		}
		if apiKey == "" {
			log.Fatal("API key cannot be empty")
		}
		cfg.GoogleAIAPIKey = apiKey
		if err := config.Save(cfg); err != nil {
			log.Printf("Warning: Failed to save config: %v", err)
		}
	}

	if cfg.AIModel == "" {
		modelName, err := config.PromptForAIModel()
		if err != nil {
			log.Fatalf("Failed to get AI Model: %v", err)
		}
		if modelName == "" {
			// Default fallback if user just presses enter?
			// Or enforce it? Let's enforce it for now as per plan, or maybe default to gemini-2.5-flash if empty?
			// The prompt says "e.g., gemini-2.5-flash", let's default to it if empty for better UX.
			modelName = "gemini-2.5-flash"
		}
		cfg.AIModel = modelName
		if err := config.Save(cfg); err != nil {
			log.Printf("Warning: Failed to save config: %v", err)
		}
	}

	// Use the configured Gemini model.
	m, err := gemini.NewModel(ctx, cfg.AIModel, &genai.ClientConfig{
		APIKey: cfg.GoogleAIAPIKey,
	})

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	return &Agent{
		llm: m,
	}
}

// Analyze analyzes the code changes.
func (a *Agent) Analyze(diff string) (string, error) {
	ctx := context.Background()

	// Construct the prompt
	prompt := fmt.Sprintf("%s\n\n%s", DefaultReviewPrompt, diff)

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Parts: []*genai.Part{
					genai.NewPartFromText(prompt),
				},
			},
		},
	}

	// Generate content
	// GenerateContent returns an iterator.
	respStream := a.llm.GenerateContent(ctx, req, false)

	var fullText string
	for resp, err := range respStream {
		if err != nil {
			return "", fmt.Errorf("failed to generate content: %w", err)
		}
		if resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part.Text != "" {
					fullText += part.Text
				}
			}
		}
	}

	if fullText == "" {
		return "", fmt.Errorf("no content generated")
	}

	return fullText, nil
}

// Fix generates fixes for the code changes.
func (a *Agent) Fix(diff string) (string, error) {
	ctx := context.Background()

	// Construct the prompt
	prompt := fmt.Sprintf("%s\n\n%s", DefaultFixPrompt, diff)

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Parts: []*genai.Part{
					genai.NewPartFromText(prompt),
				},
			},
		},
	}

	// Generate content
	respStream := a.llm.GenerateContent(ctx, req, false)

	var fullText string
	for resp, err := range respStream {
		if err != nil {
			return "", fmt.Errorf("failed to generate content: %w", err)
		}
		if resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part.Text != "" {
					fullText += part.Text
				}
			}
		}
	}

	if fullText == "" {
		return "", fmt.Errorf("no content generated")
	}

	return fullText, nil
}
