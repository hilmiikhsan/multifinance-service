package utils

import (
	"fmt"
	"time"
)

// GenerateContractNumber generates a contract number with prefix, date, customer ID, and counter
func GenerateContractNumber(customerID int) string {
	// Get current date in "YYYYMMDD" format
	date := time.Now().Format("20060102")

	// Format contract number
	contractNumber := fmt.Sprintf("TRX%s%04d", date, customerID)

	return contractNumber
}
