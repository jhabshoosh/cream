package quote
import (
	"log"
	fmp "github.com/jhabshoo/fmp/client"
)


func getFinancials(symbol string) *fmp.CompanyQuote {
	symbolInput := []string {symbol}
	quoteResponse, err := fmp.FetchCompanyQuote(symbolInput)
	if (err != nil) {
		log.Println(err)
	}
	if (quoteResponse != nil && len(quoteResponse) > 0) {
		return &quoteResponse[0]
	}
	return new(fmp.CompanyQuote)
}


func highLowerFilterRule(quote *fmp.CompanyQuote) bool {
	return (quote.YearHigh / quote.YearLow) > (1/.85)
}

func marketCapFilterRule(quote *fmp.CompanyQuote) bool {
	return quote.MarketCap >= 1000000000
}

