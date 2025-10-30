package genaistructbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	genai "google.golang.org/genai"
)

const ResponseMIMEType = "application/json"

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
		schemaJSON string,
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
func _generate[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples map[string]T,
	schema *genai.Schema,
	output *T,
) error {

	instructionsContent := &genai.Content{
		Parts: []*genai.Part{{Text: instructions}},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: instructionsContent,
		ResponseMIMEType:  ResponseMIMEType,
		ResponseSchema:    schema,
		Temperature:       float32Ptr(0.2),
	}
	parts := []*genai.Part{{Text: prompt}}
	for exPrompt, exValue := range examples {
		exampleJSON, _ := json.MarshalIndent(exValue, "", "  ")
		parts = append(parts, &genai.Part{
			Text: fmt.Sprintf("Example prompt: %s\nExpected JSON: %s", exPrompt, exampleJSON),
		})
	}

	content := []*genai.Content{{Parts: parts}}

	resp, err := llm.Models.GenerateContent(ctx, model, content, config)
	if err != nil {
		return fmt.Errorf("❌ error generating structured response: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("❌ no response received from model")
	}

	part := resp.Candidates[0].Content.Parts[0]
	raw := strings.TrimSpace(part.Text)

	if err := json.Unmarshal([]byte(raw), output); err != nil {
		return fmt.Errorf("❌ failed to unmarshal model output: %w\nRaw output: %s", err, raw)
	}

	return nil
}
func GenerateFromSchemaGeneric[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples map[string]T,
	schemaJSON string,
	output *T,
) error {
	var parsedSchema genai.Schema
	if err := json.Unmarshal([]byte(schemaJSON), &parsedSchema); err != nil {
		return fmt.Errorf("❌ failed to parse schema JSON: %w", err)
	}
	return _generate(ctx, llm, model, prompt, instructions, examples, &parsedSchema, output)
}
func (b *GenAIStructBuilder) GenerateFromchema(
	ctx context.Context,
	model string,
	prompt string,
	instructions string,
	examples map[string]map[string]interface{},
	schemaJSON string,
	output *map[string]interface{},
) error {
	return GenerateFromSchemaGeneric[map[string]interface{}](ctx, b.llm, model, prompt, instructions, examples, schemaJSON, output)
}

func GenerateFromStructGeneric[T any](
	ctx context.Context,
	llm *genai.Client,
	model string,
	prompt string,
	instructions string,
	examples map[string]T,
	output *T,
) error {
	schema := BuildSchema(output)
	return _generate(ctx, llm, model, prompt, instructions, examples, schema, output)
}
func (b *GenAIStructBuilder) GenerateFromStruct(
	ctx context.Context,
	model string,
	prompt string,
	instructions string,
	examples map[string]map[string]interface{},
	output *map[string]interface{},
) error {
	return GenerateFromStructGeneric[map[string]interface{}](ctx, b.llm, model, prompt, instructions, examples, output)
}

//
// func main() {
// 	ctx := context.Background()
//
// 	apiKey := os.Getenv("GEMINI_API_KEY")
// 	if apiKey == "" {
// 		panic("⚠️ Please set GEMINI_API_KEY environment variable")
// 	}
//
// 	client, err := genai.NewClient(ctx, &genai.ClientConfig{
// 		APIKey:  apiKey,
// 		Backend: genai.BackendGeminiAPI,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer client.Close()
//
// 	model := "gemini-2.0-flash"
//
// 	instructions := `
// You are a precise JSON field extractor for job prompts.
// Your task is to extract structured search parameters from natural language inputs about hiring or job searches.
//
// Output must strictly match this JSON schema:
// {
//   "skills": [string],
//   "location": string,
//   "job_title": string,
//   "yearsof_experience_from": int,
//   "yearsof_experience_to": int,
//   "industry": [string],
//   "company": string
// }
//
// Rules:
// - Always return valid JSON only — no markdown or explanations.
// - If information is missing, use "" for strings and [] for arrays.
// - For years of experience, infer from *seniority level* or explicit numbers as follows:
//
// Seniority Mapping:
// - "Intern" / "Entry-level" → yearsof_experience_from = 0, yearsof_experience_to = 1
// - "Junior" → 1–3
// - "Intermediate" / "Mid-level" → 3–5
// - "Senior" → 5–10
// - "Lead" / "Principal" → 8–12
// - "Manager" / "Director" → 10–15
// - If the prompt mentions explicit years (e.g., "7+ years of experience"), use that directly.
// - If both seniority and numeric years appear, prefer the numeric years.
// - If neither is provided, set both values to 0.
//
// - Location should include both city and country if provided in the prompt (e.g., "Cairo, Egypt").
// - Company should only be filled if the prompt explicitly mentions it (e.g., "at Google" → company = "Google").
// - Industry can be inferred from context words like "finance", "healthcare", "education", or "technology".
// - Skills should be extracted as specific technologies, tools, or expertise areas mentioned in the prompt.
// `
//
// 	var examples = map[string]PersonListRequest{
// 		// 🧩 Job title + seniority mapping
// 		"Looking for someone mid-level in backend work, maybe with Java and Spring.": {
// 			Skills:                []string{"Java", "Spring"},
// 			Location:              "",
// 			JobTitle:              "Backend Developer",
// 			YearsOfExperienceFrom: 3,
// 			YearsOfExperienceTo:   5,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"We need a team lead for our React project.": {
// 			Skills:                []string{"React"},
// 			Location:              "",
// 			JobTitle:              "Team Lead",
// 			YearsOfExperienceFrom: 8,
// 			YearsOfExperienceTo:   12,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Hiring junior to senior cloud engineers with AWS experience.": {
// 			Skills:                []string{"AWS", "Cloud Engineering"},
// 			Location:              "",
// 			JobTitle:              "Cloud Engineer",
// 			YearsOfExperienceFrom: 1,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Need someone who can handle accounting for a small firm.": {
// 			Skills:                []string{"Accounting"},
// 			Location:              "",
// 			JobTitle:              "Accountant",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   2,
// 			Industry:              "Finance",
// 			Company:               []string{},
// 		},
//
// 		// 🌍 Locations
// 		"Backend devs wanted in Riyadh.": {
// 			Skills:                []string{"Backend Development"},
// 			Location:              "Riyadh, Saudi Arabia",
// 			JobTitle:              "Backend Developer",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   0,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Looking for software testers in Egypt.": {
// 			Skills:                []string{"Testing"},
// 			Location:              "Egypt",
// 			JobTitle:              "Software Tester",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   0,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Hiring Flutter developers in Amman, Jordan.": {
// 			Skills:                []string{"Flutter"},
// 			Location:              "Amman, Jordan",
// 			JobTitle:              "Flutter Developer",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   0,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
//
// 		// 🏢 Companies
// 		"We want someone who worked previously at Amazon or Microsoft.": {
// 			Skills:                []string{},
// 			Location:              "",
// 			JobTitle:              "",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   0,
// 			Industry:              "",
// 			Company:               []string{"Amazon", "Microsoft"},
// 		},
// 		"Candidates from FAANG companies preferred.": {
// 			Skills:                []string{},
// 			Location:              "",
// 			JobTitle:              "",
// 			YearsOfExperienceFrom: 5,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Technology",
// 			Company:               []string{"Facebook", "Amazon", "Apple", "Netflix", "Google"},
// 		},
// 		"Senior PMs with experience at a top-tier tech giant like Meta.": {
// 			Skills:                []string{"Project Management"},
// 			Location:              "",
// 			JobTitle:              "Senior Project Manager",
// 			YearsOfExperienceFrom: 5,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Technology",
// 			Company:               []string{"Meta"},
// 		},
//
// 		// 💼 Industries
// 		"Hiring Python developers for a fintech startup.": {
// 			Skills:                []string{"Python"},
// 			Location:              "",
// 			JobTitle:              "Software Developer",
// 			YearsOfExperienceFrom: 0,
// 			YearsOfExperienceTo:   0,
// 			Industry:              "Finance",
// 			Company:               []string{},
// 		},
// 		"We’re building an AI product in healthcare.": {
// 			Skills:                []string{"AI", "Machine Learning"},
// 			Location:              "",
// 			JobTitle:              "AI Engineer",
// 			YearsOfExperienceFrom: 3,
// 			YearsOfExperienceTo:   5,
// 			Industry:              "Healthcare",
// 			Company:               []string{},
// 		},
//
// 		// 🧠 Years of experience
// 		"We need engineers with 7+ years in backend systems.": {
// 			Skills:                []string{"Backend Development"},
// 			Location:              "",
// 			JobTitle:              "Backend Engineer",
// 			YearsOfExperienceFrom: 7,
// 			YearsOfExperienceTo:   7,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Prefer people with around 3 to 5 years of experience.": {
// 			Skills:                []string{},
// 			Location:              "",
// 			JobTitle:              "",
// 			YearsOfExperienceFrom: 3,
// 			YearsOfExperienceTo:   5,
// 			Industry:              "",
// 			Company:               []string{},
// 		},
// 		"Hiring senior data scientists with 10 years experience.": {
// 			Skills:                []string{"Data Science"},
// 			Location:              "",
// 			JobTitle:              "Senior Data Scientist",
// 			YearsOfExperienceFrom: 10,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
//
// 		// 🌐 Combined prompts
// 		"Find junior Python developers in Cairo with Django skills for a fintech company.": {
// 			Skills:                []string{"Python", "Django"},
// 			Location:              "Cairo, Egypt",
// 			JobTitle:              "Junior Python Developer",
// 			YearsOfExperienceFrom: 1,
// 			YearsOfExperienceTo:   3,
// 			Industry:              "Finance",
// 			Company:               []string{},
// 		},
// 		"Senior cloud architects in Dubai familiar with AWS, GCP, and Kubernetes.": {
// 			Skills:                []string{"AWS", "GCP", "Kubernetes"},
// 			Location:              "Dubai, UAE",
// 			JobTitle:              "Senior Cloud Architect",
// 			YearsOfExperienceFrom: 5,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Technology",
// 			Company:               []string{},
// 		},
// 		"Data analysts at Amazon with 3+ years of SQL experience.": {
// 			Skills:                []string{"SQL", "Data Analysis"},
// 			Location:              "",
// 			JobTitle:              "Data Analyst",
// 			YearsOfExperienceFrom: 3,
// 			YearsOfExperienceTo:   3,
// 			Industry:              "Technology",
// 			Company:               []string{"Amazon"},
// 		},
// 		"Healthcare AI researchers with TensorFlow and NLP expertise.": {
// 			Skills:                []string{"TensorFlow", "NLP"},
// 			Location:              "",
// 			JobTitle:              "AI Researcher",
// 			YearsOfExperienceFrom: 5,
// 			YearsOfExperienceTo:   10,
// 			Industry:              "Healthcare",
// 			Company:               []string{},
// 		},
// 	}
//
// 	fmt.Println("🧠 TAL CLI — Job Prompt Extractor (Seniority-Aware)")
// 	fmt.Println("Type a job search prompt (or 'exit' to quit):\n")
//
// 	// reader := bufio.NewReader(os.Stdin)
// 	schemaJSON := `{
//   "type": "OBJECT",
//   "properties": {
//     "skills": { "type": "ARRAY", "items": { "type": "STRING" } },
//     "location": { "type": "STRING" },
//     "job_title": { "type": "STRING" },
//     "yearsof_experience_from": { "type": "INTEGER" },
//     "yearsof_experience_to": { "type": "INTEGER" },
//     "industry": { "type": "STRING" },
//     "company": { "type": "ARRAY", "items": { "type": "STRING" } }
//   },
//   "required": [
//     "skills",
//     "location",
//     "job_title",
//     "yearsof_experience_from",
//     "yearsof_experience_to",
//     "industry",
//     "company"
//   ]
// }`
//
// 	var output PersonListRequest
//
// 	err = generate_struct_from_schema(
// 		ctx,
// 		client,
// 		model,
// 		"need backend go dev 5yrs exp cairo",
// 		instructions,
// 		examples,
// 		schemaJSON,
// 		&output,
// 	)
//
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("LLM generation failed")
// 	}
//
// 	fmt.Printf("✅ Parsed result: %+v\n", output)
// 	// for {
// 	// 	fmt.Print("💬 > ")
// 	// 	prompt, _ := reader.ReadString('\n')
// 	// 	prompt = strings.TrimSpace(prompt)
// 	// 	if prompt == "" {
// 	// 		continue
// 	// 	}
// 	// 	if strings.EqualFold(prompt, "exit") {
// 	// 		fmt.Println("👋 Goodbye!")
// 	// 		break
// 	// 	}
// 	//
// 	// 	var output PersonListRequest
// 	// 	err := generate_struct(ctx, client, model, prompt, instructions, examples, &output)
// 	// 	if err != nil {
// 	// 		log.Error().Err(err).Msg("generation failed")
// 	// 		continue
// 	// 	}
// 	//
// 	// 	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
// 	// 	fmt.Println("\n🧾 Extracted JSON:")
// 	// 	fmt.Println(string(jsonBytes))
// 	// 	fmt.Println()
// 	// }
// }
//
// func float32Ptr(v float32) *float32 { return &v }
