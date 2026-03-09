package reports

import (
	"fmt"

	"context"
)

type reportServiceSqlc struct {
	reportAgent ReportAgent
	notionSvc   NotionClient
}

func NewReportService(reportAgent ReportAgent, notionSvc NotionClient) ReportService {
	return &reportServiceSqlc{reportAgent: reportAgent, notionSvc: notionSvc}
}

// GenerateReport generates a new report in the system.
func (u *reportServiceSqlc) GenerateReport(ctx context.Context, content string) (string, error) {
	// page - "31ea5909-4097-8040-86e3-c6c04293b3d9"
	// db - "31ea5909-4097-8180-ad00-f717639dafb5"
	// report object to markdown
	// dbId, err := u.notionSvc.GetOrCreateReportsDatabase(ctx, "31ea5909-4097-8040-86e3-c6c04293b3d9")
	// if err != nil {
	// 	return "", err
	// }

	// pageId, err := u.notionSvc.GetOrCreateReportsPage(ctx, dbId)
	// if err != nil {
	// 	return "", err
	// }

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

	_e := Report{
		MetaData:              reportMetaData,
		CandidateCompetencies: reportcandidateCompitencies,
		CandidateStrengths:    reportcandidateStrengths,
		CandidateConcerns:     reportcandidateConcerns,
		CandidateSignals:      reportcandidateSignals,
		OverallRecommendation: reportOverallRecommendation,
		FinalSummary:          reportFinalSummary,
	}

	fmt.Printf("%+v", _e)

	return "report-markdown", nil
}
