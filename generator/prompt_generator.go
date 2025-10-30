package generator

import (
	"context"

	"github.com/darwishdev/genaistructbuilder/internal"
	genai "google.golang.org/genai"
)

type PromptGenerator[T any] struct {
	Prompt              string
	Instructions        string
	Examples            []internal.PromptExample[T]
	CategorizedExamples map[string][]internal.PromptExample[T]
	Schema              *genai.Schema
}

func (g *PromptGenerator[T]) Execute(ctx context.Context, generateContent internal.GenerateContentFunc, model string, output *T) error {
	config := internal.GenerateConfig(ctx, g.Instructions, g.Schema)
	parts := []*genai.Part{{Text: g.Prompt}}
	internal.ExamplesHandler(parts, g.Examples, g.CategorizedExamples)
	content := []*genai.Content{{Parts: parts}}
	return internal.ExecuteLLMCall(ctx, generateContent, model, content, config, output)
}
