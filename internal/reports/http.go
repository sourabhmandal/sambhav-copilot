package reports

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ReportHandler interface {
	GenerateReport(c *gin.Context)
}

type reportHandler struct {
	reportService ReportService
}

func NewReportHandler(reportService ReportService) ReportHandler {
	return &reportHandler{reportService: reportService}
}

// GenerateReport handles HTTP requests for generating a new report.
func (h *reportHandler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// read a file content and pass it to the service layer to generate a report
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(rootDir, "example", "interview-input-1.md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	reportmd, err := h.reportService.GenerateReport(c.Request.Context(), string(data))
	if err != nil {
		panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Report generated successfully", "report": reportmd})
}
