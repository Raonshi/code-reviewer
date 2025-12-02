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
1.  **Output Language**: All explanations, summaries, analysis, and feedback must be written in **%s**.
2.  **Technical Terms**: Use original English terms for industry-standard terminology (e.g., 'Edge Case', 'Null Pointer Exception', 'Race Condition'), but provide %s explanations if the context requires clarity.

# Workflow
1.  **Analyze & Summarize**: Understand the logic changes in the provided ` + "`git diff`" + ` and provide a brief summary of the changes.
2.  **Deep Code Review**: Scrutinize the code for:
    * **Syntax & Logic**: Errors, bugs, and potential faults.
    * **Security**: Vulnerabilities (e.g., injection, data exposure).
    * **Performance**: Efficiency and resource usage.
    * **Refactoring**: Code cleanliness and maintainability.
3.  **Classify**: Evaluate each function or module based on the **Classification Criteria** below.
4.  **Propose Improvements**: For any code not rated as 'Good', provide concrete, actionable improvement guides or corrected code snippets in **%s**.

# Classification Criteria
Evaluate each change strictly according to the following levels:

1.  **Good**
    * **Definition**: Probability of errors converges to 0%%.
    * **Status**: Ready for immediate deployment. Guarantees ≥ 99%% normal operation.
2.  **Not Bad**
    * **Definition**: No immediate errors, but the code is messy or has potential risks in Edge Cases.
    * **Status**: Normal operation within intended scope. Guarantees ≥ 90%% normal operation.
3.  **Bad**
    * **Definition**: Fatal errors, bugs, security risks, or performance degradation are certain.
    * **Status**: Cannot be deployed. Guarantees < 90%% normal operation.
4.  **Need Check**
    * **Definition**: No technical errors (≥ 99%% normal operation), but business logic has significantly changed.
    * **Status**: Requires human verification to ensure it matches the planning intent.

# Output Format
Please strictly follow the format below for your report:

## [Function/Module Name]
- **Grade**: [Good / Not Bad / Bad / Need Check]
- **Summary**: (Briefly summarize the changes in this module in **%s**)
- **Analysis**: (Detailed evaluation of logic, security, and performance in **%s**)
- **Improvement Suggestions**: (Required for 'Not Bad', 'Bad', or 'Need Check'. Provide specific code fixes or refactoring advice in **%s**)

---
*(Repeat the above block for each major change)*

# Input Data
[Git Diff Data will be inserted here]
`

const DefaultFixPrompt = `
# Role
You are a **Senior Software Engineer** and **Code Review Expert**. Your goal is to provide corrected code snippets to fix issues found in the provided ` + "`git diff`" + `.

# Primary Constraints
1.  **Output Language**: The explanation should be in **%s**, but the code must be in the original language.
2.  **Scope**: Only fix the code present in the diff. Do not rewrite the entire file unless necessary.

# Workflow
1.  **Analyze**: Understand the issues in the ` + "`git diff`" + `.
2.  **Fix**: Generate the corrected code.
3.  **Explain**: Briefly explain what was fixed in **%s**.

# Output Format
Provide the fixed code in a code block, followed by a brief explanation.
`
const DefaultDocumentPrompt = `
# Role
You are a **Technical Writer** and **Software Engineer**. Your goal is to generate technical documentation for the provided ` + "`git diff`" + ` changes.

# Primary Constraints
1.  **Output Language**: All documentation must be written in **%s**.
2.  **Format**: Use Markdown.

# Workflow
1.  **Analyze**: Understand the changes in the ` + "`git diff`" + `.
2.  **Document**: Generate technical documentation explaining the changes.
    *   **Overview**: A brief summary of what changed.
    *   **Details**: Detailed explanation of the changes, including why they were made (if inferable) and how they affect the system.
    *   **Impact**: Any potential impact on other parts of the system.

# Output Format
Please strictly follow the format below:

## Overview
(Brief summary in **%s**)

## Details
(Detailed explanation in **%s**)

## Impact
(Potential impact in **%s**)
`

// Agent represents the AI agent.
type Agent struct {
	llm            model.LLM
	outputLanguage string
}

// New creates a new Agent.
func New() *Agent {
	ctx := context.Background()

	cfg := ensureConfig()
	ensureModelValidation(ctx, cfg)

	// Use the configured Gemini model.
	m, err := gemini.NewModel(ctx, cfg.AIModel, &genai.ClientConfig{
		APIKey: cfg.GoogleAIAPIKey,
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	return &Agent{
		llm:            m,
		outputLanguage: cfg.OutputLanguage,
	}
}

// ensureConfig loads the config or prompts the user for necessary information.
func ensureConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
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
	if cfg.OutputLanguage == "" {
		lang, err := config.PromptForOutputLanguage()
		if err != nil {
			log.Fatalf("Failed to get Output Language: %v", err)
		}
		if lang == "" {
			lang = "Korean"
		}
		cfg.OutputLanguage = lang
		if err := config.Save(cfg); err != nil {
			log.Printf("Warning: Failed to save config: %v", err)
		}
	}
	return cfg
}

// ensureModelValidation ensures the AI model is selected and validated.
func ensureModelValidation(ctx context.Context, cfg *config.Config) {
	if cfg.AIModel != "" {
		return
	}

	for {
		modelName, err := config.PromptForAIModel()
		if err != nil {
			log.Fatalf("Failed to get AI Model: %v", err)
		}
		if modelName == "" {
			modelName = "gemini-2.5-flash"
		}

		if validateModel(ctx, modelName, cfg.GoogleAIAPIKey) {
			fmt.Println("Model validated successfully!")
			cfg.AIModel = modelName
			if err := config.Save(cfg); err != nil {
				log.Printf("Warning: Failed to save config: %v", err)
			}
			break
		}
	}
}

// validateModel attempts to validate the model with a minimal request.
func validateModel(ctx context.Context, modelName, apiKey string) bool {
	fmt.Printf("Validating model '%s'...\n", modelName)
	tempModel, err := gemini.NewModel(ctx, modelName, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		fmt.Printf("Error creating model client: %v. Please try again.\n", err)
		return false
	}

	validateReq := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Parts: []*genai.Part{
					genai.NewPartFromText("Hi"),
				},
			},
		},
	}

	stream := tempModel.GenerateContent(ctx, validateReq, false)
	for _, err := range stream {
		if err != nil {
			fmt.Printf("Validation failed: %v. Please check the model name and try again.\n", err)
			return false
		}
	}
	return true
}

// Analyze analyzes the code changes.
func (a *Agent) Analyze(diff string) (string, error) {
	prompt := fmt.Sprintf("%s\n\n%s", fmt.Sprintf(DefaultReviewPrompt, a.outputLanguage, a.outputLanguage, a.outputLanguage, a.outputLanguage, a.outputLanguage, a.outputLanguage), diff)
	return a.generateContent(prompt)
}

// Fix generates fixes for the code changes.
func (a *Agent) Fix(diff string) (string, error) {
	prompt := fmt.Sprintf("%s\n\n%s", fmt.Sprintf(DefaultFixPrompt, a.outputLanguage, a.outputLanguage), diff)
	return a.generateContent(prompt)
}

// Document generates technical documentation for the code changes.
func (a *Agent) Document(diff string) (string, error) {
	prompt := fmt.Sprintf("%s\n\n%s", fmt.Sprintf(DefaultDocumentPrompt, a.outputLanguage, a.outputLanguage, a.outputLanguage, a.outputLanguage), diff)
	return a.generateContent(prompt)
}

// generateContent calls the LLM and aggregates the response.
func (a *Agent) generateContent(prompt string) (string, error) {
	ctx := context.Background()
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Parts: []*genai.Part{
					genai.NewPartFromText(prompt),
				},
			},
		},
	}

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
