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
