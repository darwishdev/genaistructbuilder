package genaistructbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	genai "google.golang.org/genai"
)

const ResponseMIMEType = "application/json"

type PromptExample[T any] struct {
	Prompt   string `json:"prompt"`
	Response T      `json:"response"`
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
func _generate[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples []PromptExample[T],
	categorizedExamples map[string][]PromptExample[T],
	schema *genai.Schema,
	output *T,
) error {

	instructionsContent := &genai.Content{
		Parts: []*genai.Part{{Text: instructions}},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: instructionsContent,
		ResponseMIMEType:  ResponseMIMEType,
		ResponseSchema:    schema,
		Temperature:       float32Ptr(0.2),
	}
	parts := []*genai.Part{{Text: prompt}}

	_appendExamples(&parts, examples)
	for category, exampleSlice := range categorizedExamples {
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("\n--- Categorized Example Group For Category :%d ---\n", category),
		})
		_appendExamples(&parts, exampleSlice)
	}

	content := []*genai.Content{{Parts: parts}}

	resp, err := llm.Models.GenerateContent(ctx, model, content, config)
	if err != nil {
		return fmt.Errorf("❌ error generating structured response: %w", err)
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
func GenerateFromSchemaGeneric[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples []PromptExample[T],
	categorizedExamples map[string][]PromptExample[T],
	schemaJSON []byte,
	output *T,
) error {
	var parsedSchema genai.Schema
	if err := json.Unmarshal(schemaJSON, &parsedSchema); err != nil {
		return fmt.Errorf("❌ failed to parse schema JSON: %w", err)
	}
	return _generate(ctx, llm, model, prompt, instructions, examples, categorizedExamples, &parsedSchema, output)
}
func GenerateFromStructGeneric[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples []PromptExample[T],
	categorizedExamples map[string][]PromptExample[T],
	output *T,
) error {
	schema := BuildSchema(output)
	return _generate(ctx, llm, model, prompt, instructions, examples, categorizedExamples, schema, output)
}
