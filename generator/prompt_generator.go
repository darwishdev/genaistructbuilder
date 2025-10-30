package generator

import (
	"context"

	"github.com/darwishdev/genaistructbuilder"
	"github.com/darwishdev/genaistructbuilder/internal"
	genai "google.golang.org/genai"
)

type PromptGenerator[T any] struct {
	Prompt              string
	Instructions        string
	Examples            []genaistructbuilder.PromptExample[T]
	CategorizedExamples map[string][]genaistructbuilder.PromptExample[T]
	Schema              []byte
}

func (g PromptGenerator[T]) Execute(ctx context.Context, generateContent genaistructbuilder.GenerateContentFunc, model string, output *T) error {
	schema, err := internal.BuildSchemaFromJson(g.Schema)
	if err != nil {
		return err
	}
	config := internal.GenerateConfig(ctx, g.Instructions, schema)
	parts := []*genai.Part{{Text: g.Prompt}}
	internal.ExamplesHandler(parts, g.Examples, g.CategorizedExamples)
	content := []*genai.Content{{Parts: parts}}
	return internal.ExecuteLLMCall(ctx, generateContent, model, content, config, output)
}
