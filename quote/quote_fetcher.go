package quote

import (
	"fmt"
	"github.com/jhabshoo/fmp"
)

// QuoteFetcher fetches from FMP for symbols
type QuoteFetcher struct {
	Count int
}


func getQuote(symbol string) *fmp.CompanyQuote {
	symbolInput := []string {symbol}
	quoteResponse, err := fmp.FetchCompanyQuote(symbolInput)
	if (err != nil) {
		fmt.Println(err)
	}
	if (quoteResponse != nil && len(quoteResponse) > 0) {
		return &quoteResponse[0]
	}
	return new(fmp.CompanyQuote)
}

// Process consumes from symbol channel and emits CompanyQuotes
func (qf *QuoteFetcher) Process(in <- chan string) <- chan *fmp.CompanyQuote {
	out := make(chan *fmp.CompanyQuote)
	go func() {
		for v := range in {
			quote := getQuote(v)
			qf.Count++
			out <- quote
		}
		close(out)
	}()
	return out
}