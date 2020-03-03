package quote

import (
	"github.com/jhabshoo/cream/pipeline/info"
	"fmt"
	fmp "github.com/jhabshoo/fmp/client"
)

// QuoteStage fetches from FMP for symbols
type QuoteStage struct {
	GoodCount int
	BadCount int
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
func (qf *QuoteStage) Process(in <- chan *info.Info) <- chan *fmp.CompanyQuote {
	out := make(chan *fmp.CompanyQuote)
	go func() {
		for v := range in {
			quote := getQuote(v.Symbol)
			if (highLowerFilterRule(quote) && marketCapFilterRule(quote)) {
				qf.GoodCount++
				out <- quote
			} else {
				qf.BadCount++
			}
		}
		close(out)
	}()
	return out
}