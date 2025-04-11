package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"google.golang.org/genai"
)

type GenImageReq struct {
	Instruction string `json:"instruction"`
}

// GenImageRequest sends an instruction to Gemini, saves the generated image locally,
// and returns the file path or textual response.
func GenImageRequest(ctx context.Context, instruction string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGoogleAI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	// Call the Gemini model with image modality
	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.0-flash-exp-image-generation",
		genai.Text(instruction),
		&genai.GenerateContentConfig{
			ResponseModalities: []string{"Image"},
		},
	)
	if err != nil {
		return "", fmt.Errorf("image generation error: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", errors.New("no candidates returned")
	}

	// Process each part using reflection
	for _, part := range resp.Candidates[0].Content.Parts {
		val := reflect.ValueOf(part)
		if val.Kind() == reflect.Ptr {
			elem := val.Elem()
			switch elem.Type().Name() {
			case "Blob":
				//mime := elem.FieldByName("MIMEType").String()
				data := elem.FieldByName("Data").Bytes()

				// Save image to local file
				fileExt := ".png" // could switch based on mime if needed
				outputPath := "output" + fileExt
				err := os.WriteFile(outputPath, data, 0644)
				if err != nil {
					return "", fmt.Errorf("failed to save image: %w", err)
				}

				fmt.Println("✅ Image saved to:", outputPath)
				return outputPath, nil

			case "Text":
				text := elem.FieldByName("Text").String()
				return text, nil
			}
		}
	}

	return "", errors.New("no valid image or text content returned")
}
