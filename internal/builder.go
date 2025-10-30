package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	genai "google.golang.org/genai"
)

const ResponseMIMEType = "application/json"

func ExecuteLLMCall[T any](
	ctx context.Context,
	generateContent GenerateContentFunc,
	model string,
	content []*genai.Content,
	config *genai.GenerateContentConfig,
	output *T,
) error {
	resp, err := generateContent(ctx, model, content, config)
	if err != nil {
		return fmt.Errorf("❌ error generating structured relation response: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("❌ no response received from model")
	}
	part := resp.Candidates[0].Content.Parts[0]
	raw := strings.TrimSpace(part.Text)
	if err := json.Unmarshal([]byte(raw), output); err != nil {
		return fmt.Errorf("❌ failed to unmarshal model output: %w\nRaw output: %s", err, raw)
	}
	return nil
}
func _appendRelationExamples[T any](parts *[]*genai.Part, examples []RelationExample[T]) {
	for _, example := range examples {
		exValue := example.Response
		exampleJSON, _ := json.MarshalIndent(exValue, "", "  ")

		*parts = append(*parts, &genai.Part{
			Text: fmt.Sprintf("Example Input JSON: %s\nExpected JSON: %s", example.RelationRecordJSON, exampleJSON),
		})
	}
}
func _appendExamples[T any](parts *[]*genai.Part, examples []PromptExample[T]) {
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
) *genai.GenerateContentConfig {
	instructionsContent := &genai.Content{
		Parts: []*genai.Part{{Text: instructions}},
	}
	config := &genai.GenerateContentConfig{
		SystemInstruction: instructionsContent,
		ResponseMIMEType:  ResponseMIMEType,
		ResponseSchema:    schema,
		Temperature:       float32Ptr(0.2),
	}
	return config
}

func RelationExampleHandler[T any](parts []*genai.Part, examples []RelationExample[T], categorizedExamples map[string][]RelationExample[T]) {
	_appendRelationExamples(&parts, examples)
	for category, exampleSlice := range categorizedExamples {
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("\n--- Categorized Example Group For Category :%s ---\n", category),
		})
		_appendRelationExamples(&parts, exampleSlice)
	}
}
func ExamplesHandler[T any](parts []*genai.Part, examples []PromptExample[T], categorizedExamples map[string][]PromptExample[T]) {
	_appendExamples(&parts, examples)
	for category, exampleSlice := range categorizedExamples {
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("\n--- Categorized Example Group For Category :%s ---\n", category),
		})
		_appendExamples(&parts, exampleSlice)
	}
}
