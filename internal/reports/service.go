package reports

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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
	var (
		wg         sync.WaitGroup
		workerPool chan struct{}

		mu      sync.Mutex
		res     Report
		errList []error
	)
	workerPool = make(chan struct{}, 200)

	// Metadata
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("--Generating report metadata--")
		meta, err := u.reportAgent.GenerateReportMetaData(ctx, content)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.MetaData = meta
		}
		mu.Unlock()
	}()

	// Candidate Competencies
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("--Generating report candidate competencies--")
		comp, err := u.reportAgent.GetCandidateCompetencies(ctx, content)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.CandidateCompetencies = comp
		}
		mu.Unlock()
	}()

	wg.Wait()
	if len(errList) > 0 {
		errNormalized := errors.Join(errList...)
		return "", errNormalized
	}

	// Strengths
	wg.Add(1)
	workerPool <- struct{}{}
	go func() {
		defer wg.Done()
		defer func() { <-workerPool }()
		fmt.Println("--Generating report candidate strengths--")
		strengths, err := u.reportAgent.GetCandidateStrengths(ctx, res.CandidateCompetencies)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.CandidateStrengths = strengths
		}
		mu.Unlock()
	}()

	// Concerns
	wg.Add(1)
	workerPool <- struct{}{}
	go func() {
		defer wg.Done()
		defer func() { <-workerPool }()
		fmt.Println("--Generating report candidate concerns--")
		concerns, err := u.reportAgent.GetCandidateConcerns(ctx, res.CandidateCompetencies, res.CandidateStrengths)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.CandidateConcerns = concerns
		}
		mu.Unlock()
	}()

	wg.Wait()
	if len(errList) > 0 {
		errNormalized := errors.Join(errList...)
		return "", errNormalized
	}

	// Signals
	wg.Add(1)
	workerPool <- struct{}{}
	go func() {
		defer wg.Done()
		defer func() { <-workerPool }()
		fmt.Println("--Generating report candidate signals--")
		signals, err := u.reportAgent.GetCandidateSignals(ctx, res.CandidateCompetencies, res.CandidateStrengths, res.CandidateConcerns)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.CandidateSignals = signals
		}
		mu.Unlock()
	}()

	// Overall Recommendation
	wg.Add(1)
	workerPool <- struct{}{}
	go func() {
		defer wg.Done()
		defer func() { <-workerPool }()
		fmt.Println("--Generating report overall recommendation--")
		rec, err := u.reportAgent.GetOverallRecommendation(ctx, res.CandidateCompetencies, res.CandidateStrengths, res.CandidateConcerns)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.OverallRecommendation = rec
		}
		mu.Unlock()
	}()

	// Final Summary
	wg.Add(1)
	workerPool <- struct{}{}
	go func() {
		defer wg.Done()
		defer func() { <-workerPool }()
		fmt.Println("--Generating report final summary--")
		summary, err := u.reportAgent.GetFinalSummary(ctx, res.CandidateCompetencies, res.CandidateStrengths, res.CandidateConcerns)
		mu.Lock()
		if err != nil {
			errList = append(errList, err)
		} else {
			res.FinalSummary = summary
		}
		mu.Unlock()
	}()

	wg.Wait()
	if len(errList) > 0 {
		errNormalized := errors.Join(errList...)
		return "", errNormalized
	}

	fmt.Println(res.FinalSummary)
	_e := Report{
		MetaData:              res.MetaData,
		CandidateCompetencies: res.CandidateCompetencies,
		CandidateStrengths:    res.CandidateStrengths,
		CandidateConcerns:     res.CandidateConcerns,
		CandidateSignals:      res.CandidateSignals,
		OverallRecommendation: res.OverallRecommendation,
		FinalSummary:          res.FinalSummary,
	}

	data, err := json.MarshalIndent(_e, "", "  ")
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n", data)
	// report object to markdown
	dbId, err := u.notionSvc.GetOrCreateReportsDatabase(ctx, "31ea5909-4097-8040-86e3-c6c04293b3d9")
	if err != nil {
		return "", err
	}

	pageId, err := u.notionSvc.GetOrCreateReportsPage(ctx, dbId)
	if err != nil {
		return "", err
	}

	return pageId.String(), nil
}
