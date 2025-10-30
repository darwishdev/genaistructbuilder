package genaistructbuilder

import (
	"context"

	"google.golang.org/genai"
)

type RelationExample[T any] struct {
	RelationRecordJSON string `json:"relation_record_json"` // The JSON input record (e.g., Job Offer JSON)
	Response           T      `json:"response"`             // The expected output JSON (e.g., Profile Search JSON)
}
type PromptExample[T any] struct {
	Prompt   string `json:"prompt"`
	Response T      `json:"response"`
}

type GenerateContentFunc func(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error)
type Generator[T any] interface {
	Execute(ctx context.Context, generateContent GenerateContentFunc, model string, output *T) error
}

type StructBuilderInterface[T any] interface {
	Build(generator Generator[T], model string, output *T) error
}
type GenAiStructBuilder[T any] struct {
	generateContent GenerateContentFunc
}

func NewStructBuilder[T any](generateContent GenerateContentFunc) StructBuilderInterface[T] {
	return &GenAiStructBuilder[T]{
		generateContent: generateContent,
	}
}
func (b *GenAiStructBuilder[T]) Build(generator Generator[T], model string, output *T) error {
	return generator.Execute(context.Background(), b.generateContent, model, output)
}
