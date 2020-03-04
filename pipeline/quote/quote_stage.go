package quote

import (
	"log"
	"github.com/jhabshoo/cream/pipeline/info"
	fmp "github.com/jhabshoo/fmp/client"
)

// QuoteStage fetches from FMP for symbols
type QuoteStage struct {
	GoodCount int
	BadCount int
}


func GetQuote(symbol string) *fmp.CompanyQuote {
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

// Process consumes from symbol channel and emits CompanyQuotes
func (qf *QuoteStage) Process(in <- chan *info.Info, quoteMap map[string]fmp.CompanyQuote) <- chan *fmp.CompanyQuote {
	out := make(chan *fmp.CompanyQuote)
	go func() {
		for v := range in {
			quote := GetQuote(v.Symbol)
			if (highLowerFilterRule(quote) && marketCapFilterRule(quote)) {
				qf.GoodCount++
				quoteMap[v.Symbol] = *quote
				out <- quote
			} else {
				qf.BadCount++
			}
		}
		close(out)
	}()
	return out
}