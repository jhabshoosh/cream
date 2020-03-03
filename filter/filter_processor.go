package filter

import (
	"fmt"
	"github.com/jhabshoo/fmp"
	"github.com/jhabshoo/cream/info"
)

// FundamentalFilterProcessor filters companys based on rules that can be applied to Info/Key Metrics
type FundamentalFilterProcessor struct {
	GoodCount int
	BadCount int
}

// SecondaryFilterProcessor filters companys based on rules that can be applied to Quotes
type SecondaryFilterProcessor struct {
	WellKnowCount int
	LatestEarningsCount int
	BadCount int
}

// SecondaryFilterOutputMessage is emitted by SecondaryFilterProcessor.Process
type SecondaryFilterOutputMessage struct {
	Symbol string
	WellKnown bool
}

// NewSFOutputMessage SecondaryFilterOutputMessage Constructor
func NewSFOutputMessage(symbol string, wellKnown bool) *SecondaryFilterOutputMessage {
	msg := new(SecondaryFilterOutputMessage)
	msg.Symbol = symbol
	msg.WellKnown = wellKnown
	return msg
}

func (sfom SecondaryFilterOutputMessage) String() string {
	var prefix string
	if (sfom.WellKnown) {
		prefix = "WK"
	} else {
		prefix = "LE"
	}
	return fmt.Sprintf("%s - %s", prefix, sfom.Symbol)
}

// Process consumes from an Info channel and emits passing symbols
func (ffp *FundamentalFilterProcessor) Process(in <- chan *info.Info) <- chan string {
	out := make(chan string)
	go func() {
		for v := range in {
			if (multipleFilterRule(v) && peRatioFilterRule(v)) {
				ffp.GoodCount++
				out <- v.Symbol
			} else {
				ffp.BadCount++
			}
		}
		close(out)
	}()
	return out
}

// Process consumes from a CompanyQuote channel and emits passes as SecondaryFilterOutputMessage
func (sfp *SecondaryFilterProcessor) Process(in <- chan *fmp.CompanyQuote) <- chan *SecondaryFilterOutputMessage {
	out := make(chan *SecondaryFilterOutputMessage)
	go func() {
		for v:= range in {
			if (highLowerFilterRule(v)) {
				if (marketCapFilterRule(v)) {
					sfp.WellKnowCount++
					msg := NewSFOutputMessage(v.Symbol, true)
					out <- msg
				} else {
					sfp.LatestEarningsCount++
					msg := NewSFOutputMessage(v.Symbol, false)
					out <- msg
				}
			} else {
				sfp.BadCount++
			}
		}
		close(out)
	}()
	return out
}

func multipleFilterRule(ti *info.Info) bool {
	return ti.EvOverEbitda < 10 && ti.EvOverEbitda > 2
}

func peRatioFilterRule(ti *info.Info) bool {
	return ti.PERatio < 10
}

func highLowerFilterRule(quote *fmp.CompanyQuote) bool {
	return (quote.YearHigh / quote.YearLow) > (1/.85)
}

func marketCapFilterRule(quote *fmp.CompanyQuote) bool {
	return quote.MarketCap >= 1000000000
}

