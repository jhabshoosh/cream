package quote
import fmp "github.com/jhabshoo/fmp/client"


func highLowerFilterRule(quote *fmp.CompanyQuote) bool {
	return (quote.YearHigh / quote.YearLow) > (1/.85)
}

func marketCapFilterRule(quote *fmp.CompanyQuote) bool {
	return quote.MarketCap >= 1000000000
}

