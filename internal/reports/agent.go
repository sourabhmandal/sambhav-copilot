package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type reportAgent struct {
	client *genai.Client
}

func NewReportAgent(ctx context.Context, reportAgentApiKey string) ReportAgent {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  reportAgentApiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &reportAgent{client: client}
}

func (a *reportAgent) GenerateReportMetaData(ctx context.Context, content string) (ReportMetaData, error) {
	// Call the external agent API to generate report metadata
	prompt := fmt.Sprintf("Extract report metadata from the following transcript of meeting between interviewer and candidate: %s", content)
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Schema for interview reports metadata.",
				Properties: map[string]*genai.Schema{
					"role": {
						Type:        genai.TypeString,
						Description: "The job role for which the candidate is being evaluated.",
					},
					"tags": {
						Type:        genai.TypeArray,
						Description: "List of tags related to the candidate or interview.",
						Items: &genai.Schema{
							Type: genai.TypeString,
						},
					},
					"interview_type": {
						Type:        genai.TypeString,
						Description: "The type of interview (e.g., Technical, System Design, Coding).",
					},
					"interview_date": {
						Type:        genai.TypeString,
						Description: "The date when the interview took place.",
					},
					"interviewer": {
						Type:        genai.TypeString,
						Description: "The name of the interviewer.",
					},
					"candidate": {
						Type:        genai.TypeString,
						Description: "The name of the candidate.",
					},
					"experience_claimed": {
						Type:        genai.TypeString,
						Description: "The years of experience claimed by the candidate.",
					},
				},
			},
		})
	if err != nil {
		return ReportMetaData{}, err
	}

	output := resp.Text()
	// Parse the response to extract ReportMetaData
	var meta ReportMetaData
	if err := json.Unmarshal([]byte(output), &meta); err != nil {
		return ReportMetaData{}, err
	}
	return meta, nil
}

func (a *reportAgent) GetCandidateCompetencies(ctx context.Context, content string) ([]Competency, error) {
	return []Competency{}, nil
}

func (a *reportAgent) GetCandidateStrengths(ctx context.Context, content string) ([]string, error) {
	return []string{}, nil
}

func (a *reportAgent) GetCandidateConcerns(ctx context.Context, content string) ([]string, error) {
	return []string{}, nil
}

func (a *reportAgent) GetCandidateSignals(ctx context.Context, content string) ([]Signal, error) {
	return []Signal{}, nil
}

func (a *reportAgent) GetOverallRecommendation(ctx context.Context, content string) (OverallRecommendation, error) {
	return OverallRecommendation{}, nil
}

func (a *reportAgent) GetFinalSummary(ctx context.Context, content string) (string, error) {
	return "", nil
}
