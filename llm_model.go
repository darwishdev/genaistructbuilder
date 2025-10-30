package genaistructbuilder

import (
	"context"

	genai "google.golang.org/genai"
)

type ModelService interface {
	GenerateContent(ctx context.Context, model string, contents ...*genai.Content) (*genai.GenerateContentResponse, error)
}

// FIX: Changed llm type in constructor to the abstract LLMClient interface.
func NewGenAiStructBuilder[T any](
	llm LLMClient,
) StructBuilder[T] {
	return &GenAiStructBuilder[T]{
		llm: llm,
	}
}

// Note: The concrete LLMClient would need a wrapper (or type casting) to satisfy LLMClient
// if you choose not to define LLMClient to match LLMClient methods directly.
// For now, we assume LLMClient satisfies LLMClient (meaning your internal code
// will use the interface methods).
