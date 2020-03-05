package quote

func highLowerFilterRule(quote Quote) bool {
	return (quote.YearHigh / quote.YearLow) > (1 / .85)
}

func marketCapFilterRule(quote Quote) bool {
	return quote.MarketCap >= 1000000000
}
