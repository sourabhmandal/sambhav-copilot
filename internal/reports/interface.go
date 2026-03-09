package reports

import (
	"context"
)

// Define a struct to capture JSON data from the request body
type GenerateReportRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ReportMetaData struct {
	Role              string   `json:"role"`
	Tags              []string `json:"tags"`
	InterviewType     string   `json:"interview_type"` // Technical / System Design / Coding
	InterviewDate     string   `json:"interview_date"`
	Interviewer       string   `json:"interviewer"`
	Candidate         string   `json:"candidate"`
	ExperienceClaimed string   `json:"experience_claimed"`
}

type OverallRecommendation struct {
	Verdict    string `json:"verdict"`    // Strong Hire, Hire, Lean Hire, Lean No Hire, No Hire
	Confidence string `json:"confidence"` // High, Medium, Low
	Summary    string `json:"summary"`
	RoleFit    string `json:"role_fit"` // Junior, Mid-Level, Senior, Staff
}

type Criterion struct {
	Name     string `json:"name"`
	Rating   int    `json:"rating"` // 1 to 5 - 1 No understanding, 2 Basic / weak, 3 Acceptable / average, 4 Strong, 5 Exceptional
	Evidence string `json:"evidence"`
}

type Competency struct {
	Name          string      `json:"name"`
	OverallRating string      `json:"overall_score"` // 1 to 5 - 1 No understanding, 2 Basic / weak, 3 Acceptable / average, 4 Strong, 5 Exceptional
	Criterions    []Criterion `json:"criterions"`
	Observations  []string    `json:"observations"`
}

type Signal struct {
	SignalType  string `json:"signal_type"` // Ownership, Learning Ability, Practical Experience, Communication, Collaboration
	Observation string `json:"observation"`
}

type Report struct {
	MetaData              ReportMetaData        `json:"metadata"`
	OverallRecommendation OverallRecommendation `json:"overall_recommendation"`
	CandidateCompetencies []Competency          `json:"candidate_competencies"`
	CandidateStrengths    []string              `json:"candidate_strengths"`
	CandidateConcerns     []string              `json:"candidate_concerns"`
	CandidateSignals      []Signal              `json:"candidate_signals"`
	FinalSummary          string                `json:"final_summary"`
}

type ReportService interface {
	GenerateReport(ctx context.Context, content string) (string, error)
}

type ReportAgent interface {
	GenerateReportMetaData(ctx context.Context, content string) (ReportMetaData, error)
	GetCandidateCompetencies(ctx context.Context, content string) ([]Competency, error)
	GetCandidateStrengths(ctx context.Context, competencies []Competency) ([]string, error)
	GetCandidateConcerns(ctx context.Context, competencies []Competency, strengths []string) ([]string, error)
	GetCandidateSignals(ctx context.Context, competencies []Competency, strengths []string, concerns []string) ([]Signal, error)
	GetOverallRecommendation(ctx context.Context, competencies []Competency, strengths []string, concerns []string) (OverallRecommendation, error)
	GetFinalSummary(ctx context.Context, competencies []Competency, strengths []string, concerns []string) (string, error)
}
