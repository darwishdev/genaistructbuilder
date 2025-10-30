package genaistructbuilder

import (
	"context"
	"testing"

	"github.com/darwishdev/genaistructbuilder/generator"
	genai "google.golang.org/genai"
)

// --- Shared Mock Implementations (Needed to run the tests) ---

// StructuredMockGenerateContentFunc returns predictable, valid JSON output.
func StructuredMockGenerateContentFunc(
	ctx context.Context,
	model string,
	contents []*genai.Content,
	config *genai.GenerateContentConfig,
) (*genai.GenerateContentResponse, error) {

	// The mock JSON data MUST match the JobSearchOutput struct.
	mockJSON := `{
        "skills": ["Go", "Kubernetes", "JSON-Mock"],
        "company": ["MockCorp"],
        "industry": "Testing",
        "location": "Test Bay, CA",
        "job_title": "Mock Data Engineer",
        "yearsof_experience_to": 7,
        "yearsof_experience_from": 3
    }`

	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []*genai.Part{
						{
							Text: mockJSON, // This is the payload your generator will unmarshal
						},
					},
				},
				FinishReason: genai.FinishReasonStop,
			},
		},
	}, nil
}

// CheckOutput verifies that the output struct contains the data returned by the mock.
func CheckOutput(t *testing.T, output JobSearchOutput, testName string) {
	if output.JobTitle != "Mock Data Engineer" {
		t.Errorf("❌ %s: JobTitle failed. Expected 'Mock Data Engineer', Got '%s'", testName, output.JobTitle)
	}
	if output.YearsOfExperienceFrom != 3 {
		t.Errorf("❌ %s: YearsFrom failed. Expected 3, Got %d", testName, output.YearsOfExperienceFrom)
	}
	if len(output.Skills) < 2 {
		t.Errorf("❌ %s: Skills count failed. Expected at least 2, Got %d", testName, len(output.Skills))
	}
}

// --- Main Test Function ---

func TestGenAiStructBuilder_Build_AllRealGenerators(t *testing.T) {
	testSchema := getJobSearchOutputSchema() // Get the global schema

	// 1. Setup: Inject the mock function into the builder
	// mockGenerateContent := GenerateContentFunc(StructuredMockGenerateContentFunc)
	builder := NewStructBuilder[JobSearchOutput](mockGenerateContent)
	testModel := "gemini-test-mock"

	// Define the common input for Relation and File
	inputJSON := `{"data": "some job data"}`
	inputFileBytes := []byte("PDF content")

	// --- Comprehensive Test Cases Array ---
	tests := []struct {
		Name      string
		Generator Generator[JobSearchOutput]
	}{
		{
			Name: "1. PromptGenerator_Injection_Test",
			Generator: &generator.PromptGenerator[JobSearchOutput]{
				Prompt:       "This prompt is ignored, we are testing the mock injection.",
				Instructions: "Extract fields.",
				Schema:       testSchema,
			},
		},
		{
			Name: "2. RelationGenerator_Injection_Test",
			Generator: &generator.RelationGenerator[JobSearchOutput]{
				RelationEntity:     "Job Entity",
				RelationContext:    "Map JSON to query.",
				RelationRecordJSON: inputJSON,
				Instructions:       "Extract fields.",
				Schema:             testSchema,
			},
		},
		{
			Name: "3. FileRelationGenerator_Injection_Test",
			Generator: &generator.FileRelationGenerator[JobSearchOutput]{
				RelationEntity:     "Resume File",
				RelationContext:    "Analyze file content.",
				RelationRecordFile: inputFileBytes,
				FileMIMEType:       "application/pdf",
				Instructions:       "Extract fields.",
				Schema:             testSchema,
			},
		},
	}

	// --- Execute All Test Cases ---
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			var output JobSearchOutput

			// 2. Execute the Build method
			err := builder.Build(tc.Generator, testModel, &output)

			// 3. Assertions
			if err != nil {
				t.Fatalf("❌ Build failed for %s: %v", tc.Name, err)
			}

			// Verify that the mock JSON was correctly unmarshaled by the real generator
			CheckOutput(t, output, tc.Name)
			t.Logf("✅ %s successful. Output: %+v", tc.Name, output)
		})
	}
}
