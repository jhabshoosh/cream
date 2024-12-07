package company_overview

func evOverEbitdaRule(cod CompanyOverviewData) bool {
	return cod.EvOverEbitda < 10 && cod.EvOverEbitda > 2
}

func peRatioFilterRule(cod CompanyOverviewData) bool {
	return cod.PERatio < 10
}

func highLowerFilterRule(cod CompanyOverviewData) bool {
	return (cod.YearHigh / cod.YearLow) > (1 / .85)
}

func marketCapFilterRule(cod CompanyOverviewData) bool {
	return cod.MarketCap >= 1000000000
}

func roeRule(cod CompanyOverviewData) bool {
	// val, err := strconv.ParseFloat(r.Data.ProfitabilityIndicator.ReturnOnEquity, 64)
	// if err != nil {
	// 	return false
	// }
	// return val > .15
	return cod.Roe > .15
}

func roaRule(cod CompanyOverviewData) bool {
	return cod.Roa > .05
}

func companyOverviewFilter(cod CompanyOverviewData) bool {
	return evOverEbitdaRule(cod) && peRatioFilterRule(cod) &&
		highLowerFilterRule(cod) && marketCapFilterRule(cod) &&
		roeRule(cod) && roaRule(cod)
}
