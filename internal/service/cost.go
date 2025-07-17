package service

func CalculateCost(totalSeconds float64, switchCount int) float64 {
	openingFee := 8.99
	perMinute := 8.99

	// Convert total time from seconds to minutes
	minutes := totalSeconds / 60.0

	// Apply time-based discount
	discount := 1.0
	if minutes > 60 {
		discount = 0.6 // 40% discount
	} else if minutes > 20 {
		discount = 0.8 // 20% discount
	}

	// Opening fee only for first scooter (first ride)
	totalOpeningFee := openingFee

	// Total per-minute cost with discount
	usageCost := minutes * perMinute * discount

	totalCost := totalOpeningFee + usageCost
	return totalCost
}
