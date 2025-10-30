package genaistructbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	genai "google.golang.org/genai"
)

const ResponseMIMEType = "application/json"

func _generate[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples map[string]T,
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
	for exPrompt, exValue := range examples {
		exampleJSON, _ := json.MarshalIndent(exValue, "", "  ")
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("Example prompt: %s\nExpected JSON: %s", exPrompt, exampleJSON),
		})
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
	examples map[string]T,
	schemaJSON []byte,
	output *T,
) error {
	var parsedSchema genai.Schema
	if err := json.Unmarshal(schemaJSON, &parsedSchema); err != nil {
		return fmt.Errorf("❌ failed to parse schema JSON: %w", err)
	}
	return _generate(ctx, llm, model, prompt, instructions, examples, &parsedSchema, output)
}
func GenerateFromStructGeneric[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples map[string]T,
	output *T,
) error {
	schema := BuildSchema(output)
	return _generate(ctx, llm, model, prompt, instructions, examples, schema, output)
}
