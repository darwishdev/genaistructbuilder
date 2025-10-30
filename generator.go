package genaistructbuilder

import (
	"context"

	"github.com/darwishdev/genaistructbuilder/internal"
)

type Generator[T any] interface {
	Execute(ctx context.Context, generateContent internal.GenerateContentFunc, model string, output *T) error
}

type StructBuilderInterface[T any] interface {
	Build(generator Generator[T], model string, output *T) error
}
type GenAiStructBuilder[T any] struct {
	generateContent internal.GenerateContentFunc
}

func NewStructBuilder[T any](generateContent internal.GenerateContentFunc) StructBuilderInterface[T] {
	return &GenAiStructBuilder[T]{
		generateContent: generateContent,
	}
}
func (b *GenAiStructBuilder[T]) Build(generator Generator[T], model string, output *T) error {
	return generator.Execute(context.Background(), b.generateContent, model, output)
}
