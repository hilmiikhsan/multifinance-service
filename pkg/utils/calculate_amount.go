package utils

func CalculateAdminFee(onTheRoadPrice int) int {
	percentage := 0.02 // 2% admin fee
	minFee := 50000    // Minimum admin fee

	fee := int(float64(onTheRoadPrice) * percentage)
	if fee < minFee {
		return minFee
	}

	return fee
}

func CalculateInterest(onTheRoadPrice int, tenorMonth int) int {
	interestRate := 0.01 // 1% per month
	interest := int(float64(onTheRoadPrice) * interestRate * float64(tenorMonth))
	return interest
}

func CalculateInstallment(onTheRoadPrice int, interestAmount int, tenorMonth int) int {
	totalPayable := onTheRoadPrice + interestAmount
	return totalPayable / tenorMonth
}
