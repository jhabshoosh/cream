package ratios

import (
	"strconv"
)

func roeFilter(r Ratio) bool {
	val, err := strconv.ParseFloat(r.Data.ProfitabilityIndicator.ReturnOnEquity, 64)
	if err != nil {
		return false
	}
	return val > .15
}

func roaFilter(r Ratio) bool {
	val, err := strconv.ParseFloat(r.Data.ProfitabilityIndicator.ReturnOnAssets, 64)
	if err != nil {
		return false
	}
	return val > .05
}

func deFilter(r Ratio) bool {
	val, err := strconv.ParseFloat(r.Data.Debt.DebtEquityRatio, 64)
	if err != nil {
		return false
	}
	return val <= 4.0
}
