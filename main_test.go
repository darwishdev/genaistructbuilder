package genaistructbuilder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/darwishdev/genaistructbuilder/internal"
	genai "google.golang.org/genai"
)

// --- Global Variables (Assumed to be defined in package scope) ---
func MockGenerateContentFunc(
	ctx context.Context,
	model string,
	contents []*genai.Content,
	config *genai.GenerateContentConfig,
) (*genai.GenerateContentResponse, error) {

	if model == "test-error-model" {
		return nil, errors.New("Mock API error: Test model requested failure")
	}

	// Simulate success by returning a minimal, valid response object.
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

var generateContent internal.GenerateContentFunc     // Assuming testClient is defined globally elsewhere
var mockGenerateContent internal.GenerateContentFunc // Assuming testClient is defined globally elsewhere
var testSchema *genai.Schema

const MODEL = "gemini-2.5-flash"

// --- Shared Output Struct and Schema ---

type JobSearchOutput struct {
	Skills                []string `json:"skills"`
	Company               []string `json:"company"`
	Industry              string   `json:"industry"`
	Location              string   `json:"location"`
	JobTitle              string   `json:"job_title"`
	YearsOfExperienceTo   int      `json:"yearsof_experience_to"`
	YearsOfExperienceFrom int      `json:"yearsof_experience_from"`
}

func getJobSearchOutputSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"skills":                  {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "A list of technical skills."},
			"company":                 {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "A list of specific companies."},
			"industry":                {Type: genai.TypeString, Description: "The relevant industry."},
			"location":                {Type: genai.TypeString, Description: "The geographical location."},
			"job_title":               {Type: genai.TypeString, Description: "The primary job title."},
			"yearsof_experience_to":   {Type: genai.TypeInteger, Description: "Upper bound for experience."},
			"yearsof_experience_from": {Type: genai.TypeInteger, Description: "Lower bound for experience."},
		},
		Required: []string{"job_title"},
	}
}
func TestMain(m *testing.M) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		os.Exit(m.Run())
		return
	}

	// Initialize the Global Client
	var err error
	generateContent = MockGenerateContentFunc
	mockGenerateContent = MockGenerateContentFunc

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize GenAI client in TestMain: %v", err))
	}

	// Initialize Global Schema
	testSchema = getJobSearchOutputSchema()

	// Defer closing the client after all tests run
	code := m.Run()
	// NOTE: You should ideally close the client here: testClient.Close()
	os.Exit(code)
}
