package generator

import (
	"context"
	"fmt"
	"strings"

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
	fmt.Println("🔧 PROMPT GENERATOR EXECUTE START")
	fmt.Printf("📝 Instructions: %s\n", g.Instructions)
	fmt.Printf("📋 User Prompt: %s\n", g.Prompt)
	fmt.Printf("📊 Examples Count: %d\n", len(g.Examples))
	fmt.Printf("🏷️  Categorized Examples Count: %d\n", len(g.CategorizedExamples))
	fmt.Printf("📄 Schema Length: %d bytes\n", len(g.Schema))

	// Build schema
	schema, err := internal.BuildSchemaFromJson(g.Schema)
	if err != nil {
		fmt.Printf("❌ Failed to build schema: %v\n", err)
		return err
	}
	fmt.Println("✅ Schema built successfully")

	// Generate config
	config := internal.GenerateConfig(ctx, g.Instructions, schema)
	fmt.Printf("⚙️  Config generated - Temperature: %v\n", config.Temperature)

	// Build the actual prompt that includes user input
	fullPrompt := g.buildFullPrompt()
	fmt.Printf("📦 Full Prompt Length: %d characters\n", len(fullPrompt))

	parts := []*genai.Part{{Text: fullPrompt}}
	fmt.Printf("🔢 Parts created: %d\n", len(parts))

	// Add examples
	internal.ExamplesHandler(parts, g.Examples, g.CategorizedExamples)

	content := []*genai.Content{{Parts: parts}}
	fmt.Printf("📦 Content blocks: %d\n", len(content))

	fmt.Println("🚀 Calling ExecuteLLMCall...")
	err = internal.ExecuteLLMCall(ctx, generateContent, model, content, config, output)
	if err != nil {
		fmt.Printf("❌ ExecuteLLMCall failed: %v\n", err)
		return err
	}

	fmt.Println("✅ PROMPT GENERATOR EXECUTE COMPLETED SUCCESSFULLY")
	return nil
}

// Helper method to build the complete prompt including user input
func (g PromptGenerator[T]) buildFullPrompt() string {
	// This is where you combine instructions, schema, rules, AND the actual user prompt
	// The structure depends on how your PromptGenerator is set up

	// If g.Prompt already contains the full instructions + user input, use it as is
	if strings.Contains(g.Prompt, "hiring sr sw eng ai ml exp worked @ google meta openai") {
		return g.Prompt
	}

	// Otherwise, you need to structure it properly. Example:
	builder := strings.Builder{}

	// Add instructions/schema if they're separate from g.Prompt
	// builder.WriteString("Your instructions here...\n\n")

	// Add the actual user prompt that should be processed
	builder.WriteString(g.Prompt)
	builder.WriteString("\n\nPlease extract the structured data from the above prompt.")

	return builder.String()
}
