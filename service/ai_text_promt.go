package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

type GeminiReq struct {
	Instruction   string `json:"instruction"`
	Source        string `json:"source"`
	FineGrainText string `json:"text"`
}

// GeminiReqClientPrompt interacts with the Gemini API to generate content
func GeminiReqClientPrompt(ctx context.Context, instruction string) (string, error) {
	// Ensure API key is set
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	// Initialize the Gemini API client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGoogleAI, // Assuming this is correct
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}
	fmt.Println(instruction)

	// Generate content using the Gemini model
	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.0-flash-exp",  // Ensure this model is correct
		genai.Text(instruction), // Input prompt
		nil,                     // Optional parameters, set to nil
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Check if there are any candidates
	if resp.Candidates == nil || len(resp.Candidates) == 0 {
		if resp.PromptFeedback != nil {
			return "", fmt.Errorf("no candidates generated: %v", resp.PromptFeedback)
		}
		return "", fmt.Errorf("no candidates generated and no prompt feedback available")
	}

	// Iterate through all candidates and collect their text content
	var outputParts []string
	for _, candidate := range resp.Candidates {
		if candidate.Content.Parts != nil {
			for _, part := range candidate.Content.Parts {
				outputParts = append(outputParts, part.Text)
			}
		}
	}

	// Combine all parts into a single string
	output := strings.Join(outputParts, "\n")

	return output, nil
}
