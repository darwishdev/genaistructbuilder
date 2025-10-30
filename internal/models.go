package internal

import (
	"context"

	genai "google.golang.org/genai"
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
