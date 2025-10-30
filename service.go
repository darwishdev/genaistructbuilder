package genaistructbuilder

import "context"

func (b *GenAIStructBuilder) GenerateFromchema(
	ctx context.Context,
	model string,
	prompt string,
	instructions string,
	examples []PromptExample[map[string]interface{}],
	categorizedExamples map[string][]PromptExample[map[string]interface{}],
	schemaJSON []byte,
	output *map[string]interface{},
) error {
	return GenerateFromSchemaGeneric[map[string]interface{}](ctx, b.llm, model, prompt, instructions, examples, categorizedExamples, schemaJSON, output)
}

func (b *GenAIStructBuilder) GenerateFromStruct(
	ctx context.Context,
	model string,
	prompt string,
	instructions string,
	examples []PromptExample[map[string]interface{}],
	categorizedExamples map[string][]PromptExample[map[string]interface{}],
	output *map[string]interface{},
) error {
	return GenerateFromStructGeneric[map[string]interface{}](ctx, b.llm, model, prompt, instructions, examples, categorizedExamples, output)
}
