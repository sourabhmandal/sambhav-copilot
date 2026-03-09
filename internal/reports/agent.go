package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
	prompt := fmt.Sprintf("Extract detailed list top 2 relevant candidate competencies and evaluate them based on the following transcript of meeting between interviewer and candidate: %s. Return the result as an object with a single key 'list' containing the array of competencies.", content)
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Object containing a list of candidate competencies.",
				Properties: map[string]*genai.Schema{
					"list": {
						Type:        genai.TypeArray,
						Description: "Array of candidate competencies.",
						Items: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"name": {
									Type:        genai.TypeString,
									Description: "The title of the competency. Example - Coding, System Design, Communication etc.",
									Example:     "System Design",
								},
								"overall_score": {
									Type:        genai.TypeString,
									Description: "The mean average score calculated from the 'rating' of all the competencies. Format: 1 (No understanding). Scale of 1 to 5 - 1 (No understanding), 2 (Basic / weak), 3 (Acceptable / average), 4 (Strong), 5 (Exceptional)",
									Example:     "4 (Strong)",
								},
								"criterions": {
									Type:        genai.TypeArray,
									Description: "List of criterions based on which the competency was evaluated. Example criterions for Coding competency - Code Efficiency, Code Readability, Problem Solving Approach etc.",
									Items: &genai.Schema{
										Type: genai.TypeObject,
										Properties: map[string]*genai.Schema{
											"name": {
												Type:        genai.TypeString,
												Description: "The title of the criterion. Examples - Code Efficiency, Code Readability, Problem Solving Approach etc. for Coding competency.",
												Example:     "Code Efficiency",
											},
											"rating": {
												Type:        genai.TypeNumber,
												Description: "The score given for the criterion. Scale of 1 to 5 - 1 No understanding, 2 Basic / weak, 3 Acceptable / average, 4 Strong, 5 Exceptional",
												Example:     3,
											},
											"evidence": {
												Type:        genai.TypeString,
												Description: "The specific quoted reference from the interview transcript that justifies the rating given for the criterion.",
												Example:     "The candidate solved the problem with a clear understanding by refering approach in `we have to use djkstras algorithm` showing a good understanding of basic programming concepts. However, the code they wrote had some inefficiencies and could be improved in terms of readability.",
											},
										},
									},
								},
								"observations": {
									Type:        genai.TypeArray,
									Description: "List of observations based on analysis of the transcript for the candidate.",
									Items: &genai.Schema{
										Type: genai.TypeString,
									},
								},
							},
						},
					},
				},
			},
		})
	if err != nil {
		return []Competency{}, err
	}

	output := resp.Text()
	// Parse the response to extract the list of competencies
	var wrapper struct {
		List []Competency `json:"list"`
	}
	if err := json.Unmarshal([]byte(output), &wrapper); err != nil {
		return []Competency{}, err
	}
	return wrapper.List, nil
}

func (a *reportAgent) GetCandidateStrengths(ctx context.Context, competencies []Competency) ([]string, error) {
	competenciesJson, err := json.Marshal(competencies)
	if err != nil {
		return []string{}, err
	}
	prompt := fmt.Sprintf("Extract candidate strengths from the following candidate competencies evaluation: %s", competenciesJson)
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Schema for interview reports metadata.",
				Properties: map[string]*genai.Schema{
					"list": {
						Type:        genai.TypeArray,
						Description: "List of Strengths of candidate based on competecies and scoring.",
						Items: &genai.Schema{
							Type: genai.TypeString,
						},
					},
				},
			},
		})
	if err != nil {
		return []string{}, err
	}

	output := resp.Text()
	// Parse the response to extract ReportMetaData
	var strengthWrapper struct {
		List []string `json:"list"`
	}
	if err := json.Unmarshal([]byte(output), &strengthWrapper); err != nil {
		return []string{}, err
	}
	return strengthWrapper.List, nil
}

func (a *reportAgent) GetCandidateConcerns(ctx context.Context, competencies []Competency, strengths []string) ([]string, error) {
	competenciesJson, err := json.Marshal(competencies)
	if err != nil {
		return []string{}, err
	}
	prompt := fmt.Sprintf("Extract candidate concerns from the following candidate competencies evaluation: %s, which is not mentioned in candidate strengths below: %s", competenciesJson, strings.Join(strengths, ", "))
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Schema for interview reports metadata.",
				Properties: map[string]*genai.Schema{
					"list": {
						Type:        genai.TypeArray,
						Description: "List of Strengths of candidate based on competecies and scoring.",
						Items: &genai.Schema{
							Type: genai.TypeString,
						},
					},
				},
			},
		})
	if err != nil {
		return []string{}, err
	}

	output := resp.Text()
	// Parse the response to extract ReportMetaData
	var strengthWrapper struct {
		List []string `json:"list"`
	}
	if err := json.Unmarshal([]byte(output), &strengthWrapper); err != nil {
		return []string{}, err
	}
	return strengthWrapper.List, nil
}

func (a *reportAgent) GetCandidateSignals(ctx context.Context, competencies []Competency, strengths []string, concerns []string) ([]Signal, error) {
	competenciesJson, err := json.Marshal(competencies)
	if err != nil {
		return []Signal{}, err
	}
	prompt := fmt.Sprintf("For these candidate signals: Ownership, Learning Ability, Practical Experience, Communication, Collaboration, give observations for each of them based on from the following candidate competencies evaluation: %s\nstrengths: %s\nweaknesses: %s", competenciesJson, strings.Join(strengths, ", "), strings.Join(concerns, ", "))
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Schema for interview reports metadata.",
				Properties: map[string]*genai.Schema{
					"list": {
						Type:        genai.TypeArray,
						Description: "List of Signals of candidate based on competecies, scoring, strengths and concerns.",
						Items: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"signal_type": {
									Type:        genai.TypeString,
									Description: "The type of signal observed. Examples include Ownership, Learning Ability, Practical Experience, Communication, Collaboration etc.",
									Example:     "Ownership",
									Enum: []string{
										"Ownership",
										"Learning Ability",
										"Practical Experience",
										"Communication",
										"Collaboration",
									},
								},
								"observation": {
									Type:        genai.TypeString,
									Description: "The specific observation from the interview transcript that justifies the signal.",
									Example:     "The candidate took ownership of a problem during the system design discussion by proactively identifying potential bottlenecks and suggesting improvements without being prompted by the interviewer.",
								},
							},
						},
					},
				},
			},
		})
	if err != nil {
		return []Signal{}, err
	}

	output := resp.Text()
	// Parse the response to extract ReportMetaData
	var signalWrapper struct {
		List []Signal `json:"list"`
	}
	if err := json.Unmarshal([]byte(output), &signalWrapper); err != nil {
		return []Signal{}, err
	}
	return signalWrapper.List, nil
}

func (a *reportAgent) GetOverallRecommendation(ctx context.Context, competencies []Competency, strengths []string, concerns []string) (OverallRecommendation, error) {
	competenciesJson, err := json.Marshal(competencies)
	if err != nil {
		return OverallRecommendation{}, err
	}
	prompt := fmt.Sprintf("Generate an overall recommendation based on the following candidate competencies evaluation: %s\nstrengths: %s\nweaknesses: %s", competenciesJson, strings.Join(strengths, ", "), strings.Join(concerns, ", "))
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeObject,
				Description: "Schema for interview reports metadata.",
				Properties: map[string]*genai.Schema{
					"list": {
						Type:        genai.TypeArray,
						Description: "List of Signals of candidate based on competecies, scoring, strengths and concerns.",
						Items: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"verdict": {
									Type:        genai.TypeString,
									Description: "The type of signal observed. Examples include Ownership, Learning Ability, Practical Experience, Communication, Collaboration etc.",
									Example:     "Ownership",
									Enum: []string{
										"Strong Hire",
										"Hire",
										"Lean Hire",
										"Lean No Hire",
										"No Hire",
									},
								},
								"confidence": {
									Type:        genai.TypeString,
									Description: "The specific observation from the interview transcript that justifies the signal.",
									Enum: []string{
										"High",
										"Medium",
										"Low",
									},
									Example: "The candidate took ownership of a problem during the system design discussion by proactively identifying potential bottlenecks and suggesting improvements without being prompted by the interviewer.",
								},
								"summary": {
									Type:        genai.TypeString,
									Description: "The specific observation from the interview transcript that justifies the signal.",
									Example:     "The candidate took ownership of a problem during the system design discussion by proactively identifying potential bottlenecks and suggesting improvements without being prompted by the interviewer.",
								},
								"role_fit": {
									Type:        genai.TypeString,
									Description: "The specific observation from the interview transcript that justifies the signal.",
									Example:     "Mid-Level",
									Enum: []string{
										"Entry-Level",
										"Mid-Level",
										"Senior-Level",
										"Staff-Level",
										"Principal-Level",
										"Distinguished-Level",
									},
								},
							},
						},
					},
				},
			},
		})
	if err != nil {
		return OverallRecommendation{}, err
	}

	output := resp.Text()
	var overallRecommendation OverallRecommendation
	if err := json.Unmarshal([]byte(output), &overallRecommendation); err != nil {
		return OverallRecommendation{}, err
	}
	return overallRecommendation, nil
}

func (a *reportAgent) GetFinalSummary(ctx context.Context, competencies []Competency, strengths []string, concerns []string) (string, error) {
	competenciesJson, err := json.Marshal(competencies)
	if err != nil {
		return "", err
	}
	prompt := fmt.Sprintf("Generate an Final Summary of candidate interview performance based on the following candidate competencies evaluation: %s\nstrengths: %s\nweaknesses: %s", competenciesJson, strings.Join(strengths, ", "), strings.Join(concerns, ", "))
	resp, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type:        genai.TypeString,
				Description: "Final Summary of the candidate's interview performance.",
			},
		})
	if err != nil {
		return "", err
	}

	output := resp.Text()
	output = strings.Trim(output, "\"") // Remove leading and trailing quotes if present
	return output, nil
}
