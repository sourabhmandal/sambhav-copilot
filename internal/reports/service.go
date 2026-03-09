package reports

import (
	"encoding/json"
	"fmt"

	"context"
)

type reportServiceSqlc struct {
	reportAgent ReportAgent
}

func NewReportService(reportAgent ReportAgent) ReportService {
	return &reportServiceSqlc{reportAgent: reportAgent}
}

// GenerateReport generates a new report in the system.
func (u *reportServiceSqlc) GenerateReport(ctx context.Context, content string) (string, error) {
	fmt.Println("--Generating report metadata--")
	reportMetaData, err := u.reportAgent.GenerateReportMetaData(ctx, content)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report candidate competencies--")
	reportcandidateCompitencies, err := u.reportAgent.GetCandidateCompetencies(ctx, content)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report candidate strengths--")
	reportcandidateStrengths, err := u.reportAgent.GetCandidateStrengths(ctx, reportcandidateCompitencies)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report candidate concerns--")
	reportcandidateConcerns, err := u.reportAgent.GetCandidateConcerns(ctx, reportcandidateCompitencies, reportcandidateStrengths)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report candidate signals--")
	reportcandidateSignals, err := u.reportAgent.GetCandidateSignals(ctx, reportcandidateCompitencies, reportcandidateStrengths, reportcandidateConcerns)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report overall recommendation--")
	reportOverallRecommendation, err := u.reportAgent.GetOverallRecommendation(ctx, reportcandidateCompitencies, reportcandidateStrengths, reportcandidateConcerns)
	if err != nil {
		return "", err
	}
	fmt.Println("--Generating report final summary--")
	reportFinalSummary, err := u.reportAgent.GetFinalSummary(ctx, reportcandidateCompitencies, reportcandidateStrengths, reportcandidateConcerns)
	if err != nil {
		return "", err
	}
	fmt.Println(reportFinalSummary)

	report := Report{
		MetaData:              reportMetaData,
		CandidateCompetencies: reportcandidateCompitencies,
		CandidateStrengths:    reportcandidateStrengths,
		CandidateConcerns:     reportcandidateConcerns,
		CandidateSignals:      reportcandidateSignals,
		OverallRecommendation: reportOverallRecommendation,
		FinalSummary:          reportFinalSummary,
	}

	reportJson, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	fmt.Println(string(reportJson))

	return "report-markdown", nil
}
