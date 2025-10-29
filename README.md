# `genaistructbuilder`

A Go package designed to simplify the generation of structured JSON output from Google's Gemini models using the `google.golang.org/genai` SDK.

This utility provides functions to:

1.  Generate a structured JSON response based on a **JSON Schema string**.
2.  Generate a structured JSON response based on a **Go struct** (by internally deriving the JSON Schema).

## ✨ Features

- **Structured Output:** Guarantees the model's response adheres to a specified JSON schema.
- **Schema Flexibility:** Supports defining the structure via raw JSON Schema or a Go struct.
- **Few-Shot Learning:** Easily incorporate examples of prompts and expected JSON outputs to guide the model.
- **Generic Functions:** Provides generic functions (`GenerateFromSchemaGeneric`, `GenerateFromStructGeneric`) for type safety when working with custom Go structs.

## 🚀 Installation

Since this package relies on the `google.golang.org/genai` SDK, ensure you have a standard Go environment.

```bash
go get github.com/darwishdev/genaistructbuilder
```

## 🛠️ Usage

### Initialization

First, you need a `genai.Client` initialized, typically by setting up your API key.

```go
package main

import (
	"context"
	"log"

	"github.com/darwishdev/genaistructbuilder"
	genai "google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// Initialize the GenAI client.
	// Ensure the GEMINI_API_KEY environment variable is set.
	client, err := genai.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Initialize the Struct Builder
	builder := genaistructbuilder.NewGenAIStructBuilder(client)

    // ... use builder functions below
}
```

### Example 1: Generate from JSON Schema

This is useful when you have the schema as a string or a file.

```go
// ... inside main or a function with `builder` and `client`

func ExampleGenerateFromSchema(ctx context.Context, builder genaistructbuilder.GenAIStructBuilderInterface, client *genai.Client) {
    const model = "gemini-2.5-flash"
    const prompt = "Please create a detailed product specification for a new smart toothbrush that uses AI."
    const instructions = "You are a Product Specification Generator. Your only job is to generate a JSON object."

    // Define the desired output structure as a JSON Schema string
    schemaJSON := `{
        "type": "object",
        "properties": {
            "productName": {"type": "string", "description": "The name of the product."},
            "features": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "name": {"type": "string"},
                        "description": {"type": "string"}
                    },
                    "required": ["name", "description"]
                }
            },
            "priceEstimate": {"type": "number", "description": "Estimated retail price in USD."}
        },
        "required": ["productName", "features"]
    }`

    // Few-shot examples (optional)
    examples := map[string]map[string]interface{}{
        "a new budget fitness tracker": {
            "productName": "PulseLite HR",
            "features": []interface{}{
                map[string]interface{}{"name": "Heart Rate Monitoring", "description": "24/7 continuous heart rate tracking."},
            },
            "priceEstimate": 49.99,
        },
    }

    // The output will be unmarshaled into this map
    var output map[string]interface{}

    log.Println("Generating product spec from schema...")
    err := builder.GenerateFromchema(
        ctx,
        client,
        model,
        prompt,
        instructions,
        examples,
        schemaJSON,
        &output,
    )
    if err != nil {
        log.Fatalf("❌ Generation failed: %v", err)
    }

    log.Printf("✅ Generated Product Spec:\n%+v\n", output)
}
```

### Example 2: Generate from Go Struct (Generic Function)

This is the preferred method for type safety, using the generic function `GenerateFromStructGeneric`.

```go
// ... inside main or a function with `builder` and `client`

// 1. Define the target Go struct
type SmartToothbrushSpec struct {
    ProductName   string    `json:"productName"`
    ModelNumber   string    `json:"modelNumber"`
    KeyFeatures   []Feature `json:"keyFeatures"`
    PriceEstimate float64   `json:"priceEstimate"`
}

type Feature struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

func ExampleGenerateFromStruct(ctx context.Context, client *genai.Client) {
    const model = "gemini-2.5-flash"
    const prompt = "Create a specification for a high-end, AI-powered electric toothbrush."
    const instructions = "You must strictly output a JSON object that matches the provided Go struct definition."

    // Few-shot examples (optional)
    examples := map[string]SmartToothbrushSpec{
        "a simple manual toothbrush": {
            ProductName: "SimpleClean 100",
            ModelNumber: "SC-100",
            KeyFeatures: []Feature{{Name: "Ergonomic Grip", Description: "Comfortable, non-slip handle."}},
            PriceEstimate: 5.99,
        },
    }

    // The output will be unmarshaled into this struct
    var output SmartToothbrushSpec

    log.Println("Generating toothbrush spec from Go struct...")
    err := genaistructbuilder.GenerateFromStructGeneric(
        ctx,
        client,
        model,
        prompt,
        instructions,
        examples,
        &output, // Pass pointer to the struct
    )
    if err != nil {
        log.Fatalf("❌ Struct generation failed: %v", err)
    }

    log.Printf("✅ Generated Smart Toothbrush Spec (Type-Safe):\n")
    log.Printf("Product: %s (Model: %s)\n", output.ProductName, output.ModelNumber)
    for _, f := range output.KeyFeatures {
        log.Printf("  - %s: %s\n", f.Name, f.Description)
    }
    log.Printf("Estimated Price: $%.2f\n", output.PriceEstimate)
}
```

---

Would you like me to elaborate on any specific part of the code, such as the internal logic of `BuildSchema(output)` which is required for the struct-based generation?
