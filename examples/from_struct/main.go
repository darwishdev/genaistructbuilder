package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/darwishdev/genaistructbuilder"
	genai "google.golang.org/genai"
)

// The model we will use for generation
const MODEL = "gemini-2.5-flash"

// --- 1. Define the Target Go Structs ---

// EmployeeInfo is the root struct for the desired JSON output.
type EmployeeInfo struct {
	EmployeeID string    `json:"employeeId"`
	FullName   string    `json:"fullName"`
	Position   string    `json:"position" description:"The official job title."`
	Skills     []string  `json:"skills" description:"A list of 3-5 technical and soft skills."`
	Projects   []Project `json:"projects"`
	IsManager  bool      `json:"isManager"`
}

// Project defines a nested object within the EmployeeInfo struct.
type Project struct {
	Name        string `json:"name"`
	DurationMos int    `json:"durationMonths" description:"Duration in months."`
}

func main() {
	// 2. Initialize GenAI Client
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("⚠️ Please set GEMINI_API_KEY environment variable")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	// --- 3. Define Generation Parameters ---

	prompt := "Create a detailed profile for a Senior Software Engineer named Alex who works primarily on Go and Kubernetes, and manages a small team."
	instructions := "You are an HR Profile Generator. Fill in the required fields to create a complete employee profile."

	// Few-shot examples (optional, using the type-safe struct)
	examples := map[string]EmployeeInfo{
		"Profile for a Junior UX Designer named Sarah who just finished one project.": {
			EmployeeID: "UX-102",
			FullName:   "Sarah Connor",
			Position:   "Junior UX Designer",
			Skills:     []string{"Figma", "User Research", "Prototyping"},
			Projects:   []Project{{Name: "Mobile App Redesign", DurationMos: 4}},
			IsManager:  false,
		},
	}

	// 4. Call the Generic Builder Function
	var output EmployeeInfo // The output is a type-safe struct

	fmt.Println("➡️ Requesting structured employee profile from Go Struct...")

	err = genaistructbuilder.GenerateFromStructGeneric(
		ctx,
		client,
		MODEL,
		prompt,
		instructions,
		examples,
		&output, // Pass a pointer to the type-safe struct
	)

	// 5. Handle Response
	if err != nil {
		fmt.Printf("\n❌ Error generating structured response: %v\n", err)
		os.Exit(1)
	}

	// Print the result from the type-safe struct
	fmt.Printf("\n✅ Successfully Generated Employee Profile:\n")
	fmt.Printf("ID: %s\n", output.EmployeeID)
	fmt.Printf("Name: %s\n", output.FullName)
	fmt.Printf("Position: %s\n", output.Position)
	fmt.Printf("Is Manager: %t\n", output.IsManager)

	// You can also marshal it to JSON to show the final output
	outputJSON, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("--- JSON Output ---")
	fmt.Println(string(outputJSON))
}
