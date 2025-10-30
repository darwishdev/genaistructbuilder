package generator

import (
	"context"
	"fmt"

	"github.com/darwishdev/genaistructbuilder/internal"
	genai "google.golang.org/genai"
)

type RelationGenerator[T any] struct {
	RelationEntity      string
	RelationContext     string
	RelationRecordJSON  string
	Instructions        string
	Examples            []internal.RelationExample[T]
	CategorizedExamples map[string][]internal.RelationExample[T]
	Schema              []byte
}

func (g *RelationGenerator[T]) Execute(ctx context.Context, generateContent internal.GenerateContentFunc, model string, output *T) error {

	schema, err := internal.BuildSchemaFromJson(g.Schema)
	if err != nil {
		return err
	}
	config := internal.GenerateConfig(ctx, g.Instructions, schema)
	mainPrompt := fmt.Sprintf(
		"Task: Generate a %s record based on the provided input JSON.\nContext: %s\nInput JSON: %s",
		g.RelationEntity,
		g.RelationContext,
		g.RelationRecordJSON,
	)
	parts := []*genai.Part{{Text: mainPrompt}}
	internal.RelationExampleHandler(parts, g.Examples, g.CategorizedExamples)
	content := []*genai.Content{{Parts: parts}}
	return internal.ExecuteLLMCall(ctx, generateContent, model, content, config, output)
}
