package genaistructbuilder

import (
	"context"

	genai "google.golang.org/genai"
)

type GenAIStructBuilder struct {
	llm *genai.Client
}

type GenAIStructBuilderInterface interface {
	GenerateFromchema(
		ctx context.Context,
		model string,
		prompt string,
		instructions string,
		examples map[string]map[string]interface{},
		schemaJSON []byte,
		output *map[string]interface{},
	) error
	GenerateFromStruct(
		ctx context.Context,
		model string,
		prompt string,
		instructions string,
		examples map[string]map[string]interface{},
		output *map[string]interface{},
	) error
}

func NewGenAIStructBuilder(llm *genai.Client) GenAIStructBuilderInterface {
	return &GenAIStructBuilder{
		llm: llm,
	}

}
