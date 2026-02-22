package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EligibilityRequest struct {
	CRNumber       string  `json:"cr_number"`
	MonthlyRevenue float64 `json:"monthly_revenue"`
}

type EligibilityResponse struct {
	Eligible      bool    `json:"eligible"`
	MaxLoanAmount float64 `json:"max_loan_amount"`
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/api/v1/eligibility", func(c *gin.Context) {
		var req EligibilityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logError(fmt.Sprintf("Failed to parse request: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logInfo(fmt.Sprintf("Received eligibility request: cr_number=%s, monthly_revenue=%.2f", req.CRNumber, req.MonthlyRevenue))

		eligible := true

		// Check if the business meets minimum revenue threshold
		if req.MonthlyRevenue < 10000 {
			eligible = false
		}

		var maxLoanAmount float64
		if eligible {
			maxLoanAmount = req.MonthlyRevenue * 3
		}

		logInfo(fmt.Sprintf("Eligibility result: cr_number=%s, eligible=%v, max_loan_amount=%.2f", req.CRNumber, eligible, maxLoanAmount))

		c.JSON(http.StatusOK, EligibilityResponse{
			Eligible:      eligible,
			MaxLoanAmount: maxLoanAmount,
		})
	})

	logInfo("loan_service starting on :3002")
	r.Run(":3002")
}
