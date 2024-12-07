package debt_equity

func debtEquityFilter(de DebtEquityData) bool {
	return (de.TotalLiabilities / de.TotalShareholderEquity) <= 4
}
