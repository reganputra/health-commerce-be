package handlers

import (
	"net/http"
	"strconv"

	"health-store/service"

	"github.com/gin-gonic/gin"
)

// GenerateReport handles report generation with various options
func GenerateReport(reportService *service.ReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		reportType := c.DefaultQuery("type", "summary")      // summary, detailed, financial
		format := c.DefaultQuery("format", "pdf")            // pdf, csv
		startDate := c.Query("start_date")                   // Optional: YYYY-MM-DD
		endDate := c.Query("end_date")                       // Optional: YYYY-MM-DD
		limitStr := c.DefaultQuery("limit", "10")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}

		// Build report request
		req := service.ReportRequest{
			ReportType: reportType,
			Format:     format,
			Limit:      limit,
		}

		if startDate != "" {
			req.StartDate = &startDate
		}
		if endDate != "" {
			req.EndDate = &endDate
		}

		// Generate report
		reportData, err := reportService.GenerateReport(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report: " + err.Error()})
			return
		}

		// Set appropriate headers based on format
		switch format {
		case "csv":
			c.Writer.Header().Set("Content-Type", "text/csv")
			c.Writer.Header().Set("Content-Disposition", "attachment; filename=report.csv")
		case "pdf":
			fallthrough
		default:
			c.Writer.Header().Set("Content-Type", "application/pdf")
			c.Writer.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
		}

		// Write report to response
		c.Writer.Write(reportData)
	}
}

// GenerateTransactionReport handles simple report generation (backward compatibility)
func GenerateTransactionReport(reportService *service.ReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate PDF through service
		pdfData, err := reportService.GenerateTransactionReport()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
			return
		}

		// Set headers for PDF response
		c.Writer.Header().Set("Content-Type", "application/pdf")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=report.pdf")

		// Write PDF to response
		c.Writer.Write(pdfData)
	}
}
