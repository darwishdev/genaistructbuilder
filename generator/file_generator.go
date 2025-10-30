package generator

import (
	"context"
	"fmt"

	"github.com/darwishdev/genaistructbuilder/internal"
	genai "google.golang.org/genai"
)

type FileRelationGenerator[T any] struct {
	RelationEntity      string
	RelationContext     string
	RelationRecordFile  []byte
	FileMIMEType        string
	Instructions        string
	Examples            []internal.RelationExample[T]
	CategorizedExamples map[string][]internal.RelationExample[T]
	Schema              []byte
}

func (g *FileRelationGenerator[T]) Execute(ctx context.Context, generateContent internal.GenerateContentFunc, model string, output *T) error {
	genSchema, err := internal.BuildSchemaFromJson(g.Schema)
	if err != nil {
		return err
	}
	config := internal.GenerateConfig(ctx, g.Instructions, genSchema)
	processedText, mediaPart, err := internal.FileAdapter(ctx, g.RelationRecordFile, g.FileMIMEType)
	if err != nil {
		return fmt.Errorf("❌ file adapter failed: %w", err)
	}
	mainPromptText := fmt.Sprintf(
		"Task: Generate a %s record based on the provided file content. \nContext: %s",
		g.RelationEntity,
		g.RelationContext,
	)
	parts := []*genai.Part{{Text: mainPromptText}}
	if mediaPart != nil {
		parts = append(parts, mediaPart)
		parts = append(parts, &genai.Part{Text: "\nInput File: (See attached media part)"})
	} else {
		parts = append(parts, &genai.Part{Text: fmt.Sprintf("\nInput File Content:\n%s", processedText)})
	}
	internal.RelationExampleHandler(parts, g.Examples, g.CategorizedExamples)
	content := []*genai.Content{{Parts: parts}}
	return internal.ExecuteLLMCall(ctx, generateContent, model, content, config, output)
}
