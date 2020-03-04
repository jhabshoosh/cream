package ratios

import (
	"log"
	fmp "github.com/jhabshoo/fmp/client"
)

type Ratio struct {
	Data fmp.FinancialRatio
	Symbol string
}

func NewRatioFromResponse(r fmp.FinancialRatiosResponse) *Ratio {
	ratio := new(Ratio)
	ratio.Symbol = r.Symbol
	if (len(r.Ratios) > 0) {
		ratio.Data = r.Ratios[0]
	}
	return ratio
}

type RatiosStage struct {
	GoodCount int
	BadCount int
}

func getRatio(symbol string) *Ratio {
	frResponse, err := fmp.FetchFinancialRatios(symbol)
	if (err != nil) {
		log.Println(err)
	}
	return NewRatioFromResponse(frResponse)
}


// Process consumes from symbol channel and emits CompanyQuotes
func (rs *RatiosStage) Process(in <- chan *fmp.CompanyQuote) <- chan string {
	out := make(chan string)
	go func() {
		for v := range in {
			ratio := getRatio(v.Symbol)
			if (roeFilter(ratio) && roaFilter(ratio)) {
				rs.GoodCount++
				out <- ratio.Symbol
			} else {
				rs.BadCount++
			}
		}
		close(out)
	}()
	return out
}