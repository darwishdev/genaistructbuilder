package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/darwishdev/genaistructbuilder"
	genai "google.golang.org/genai"
)

const ResponseMIMEType = "application/json"

func ExecuteLLMCall[T any](
	ctx context.Context,
	generateContent genaistructbuilder.GenerateContentFunc,
	model string,
	content []*genai.Content,
	config *genai.GenerateContentConfig,
	output *T,
) error {
	fmt.Println("=== LLM REQUEST CONTENT ===")
	fmt.Printf("Model: %s\n", model)

	// Print content with better formatting
	for i, c := range content {
		fmt.Printf("Content [%d]:\n", i)
		for j, part := range c.Parts {
			if part.Text != "" {
				fmt.Printf("  Part [%d] (Text):\n", j)
				fmt.Printf("    %s\n", part.Text)
			}
		}
	}

	// Print config as JSON for better readability
	if config != nil {
		configJSON, _ := json.MarshalIndent(config, "  ", "  ")
		fmt.Printf("Generation Config:\n  %s\n", string(configJSON))
	}
	fmt.Println("===========================")
	resp, err := generateContent(ctx, model, content, config)
	if err != nil {
		return fmt.Errorf("❌ error generating structured relation response: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("❌ no response received from model")
	}
	part := resp.Candidates[0].Content.Parts[0]
	raw := strings.TrimSpace(part.Text)

	fmt.Println("Raw Response")
	fmt.Println(raw)
	if err := json.Unmarshal([]byte(raw), output); err != nil {
		return fmt.Errorf("❌ failed to unmarshal model output: %w\nRaw output: %s", err, raw)
	}
	fmt.Println("Raw Response")
	fmt.Println(raw)
	return nil
}
func _appendRelationExamples[T any](parts *[]*genai.Part, examples []genaistructbuilder.RelationExample[T]) {
	for _, example := range examples {
		exValue := example.Response
		exampleJSON, _ := json.MarshalIndent(exValue, "", "  ")

		*parts = append(*parts, &genai.Part{
			Text: fmt.Sprintf("Example Input JSON: %s\nExpected JSON: %s", example.RelationRecordJSON, exampleJSON),
		})
	}
}
func _appendExamples[T any](parts *[]*genai.Part, examples []genaistructbuilder.PromptExample[T]) {
	for _, example := range examples {
		exValue := example.Response
		exampleJSON, _ := json.MarshalIndent(exValue, "", "  ")

		*parts = append(*parts, &genai.Part{
			Text: fmt.Sprintf("Example prompt: %s\nExpected JSON: %s", example.Prompt, exampleJSON),
		})
	}
}

func GenerateConfig(
	ctx context.Context,
	instructions string,
	schema *genai.Schema,
	temprature float32,
) *genai.GenerateContentConfig {
	instructionsContent := &genai.Content{
		Parts: []*genai.Part{{Text: instructions}},
	}
	config := &genai.GenerateContentConfig{
		SystemInstruction: instructionsContent,
		ResponseMIMEType:  ResponseMIMEType,
		ResponseSchema:    schema,
		Temperature:       float32Ptr(temprature),
	}
	return config
}

func RelationExampleHandler[T any](parts []*genai.Part, examples []genaistructbuilder.RelationExample[T], categorizedExamples map[string][]genaistructbuilder.RelationExample[T]) {
	_appendRelationExamples(&parts, examples)
	for category, exampleSlice := range categorizedExamples {
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("\n--- Categorized Example Group For Category :%s ---\n", category),
		})
		_appendRelationExamples(&parts, exampleSlice)
	}
}
func ExamplesHandler[T any](parts []*genai.Part, examples []genaistructbuilder.PromptExample[T], categorizedExamples map[string][]genaistructbuilder.PromptExample[T]) {
	_appendExamples(&parts, examples)
	for category, exampleSlice := range categorizedExamples {
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("\n--- Categorized Example Group For Category :%s ---\n", category),
		})
		_appendExamples(&parts, exampleSlice)
	}
}
