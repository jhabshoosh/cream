package main

import (
	"github.com/jhabshoo/cream/fmp"
	"fmt"
	"time"
	"log"
	s "strings"
	"github.com/jhabshoo/cream/ticker"
)

var tickersToProcess = make(chan string, 1000)
var processCount int = 0

func gen(tickers []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, v := range tickers {
			out <- v
		}
		close(out)
	}()
	return out
}

func main() {
	start := time.Now()

	companies := fmp.GetSymbolsList()
	var tickers []string
	for _,v := range companies {
		if (!s.Contains(v.Symbol, ".")) {
			tickers = append(tickers, v.Symbol)
		}
	}
	
	infoFetcher := new(ticker.InfoFetcher)
	tickerChan := gen(tickers)
	ic1 := infoFetcher.Process(tickerChan)
	ic2 := infoFetcher.Process(tickerChan)
	ic3 := infoFetcher.Process(tickerChan)
	ic4 := infoFetcher.Process(tickerChan)
	ic5 := infoFetcher.Process(tickerChan)
	ic6 := infoFetcher.Process(tickerChan)
	ic7 := infoFetcher.Process(tickerChan)
	ic8 := infoFetcher.Process(tickerChan)
	ic9 := infoFetcher.Process(tickerChan)
	ic10 := infoFetcher.Process(tickerChan)
	mergedTickerChan := ticker.MergeInfoChannel(ic1, ic2, ic3, ic4, ic5, ic6, ic7, ic8, ic9, ic10)


	ffp := new(ticker.FundamentalFilterProcessor)
	ffc1 := ffp.Process(mergedTickerChan)
	ffc2 := ffp.Process(mergedTickerChan)
	ffc3 := ffp.Process(mergedTickerChan)
	mergedFFChan := ticker.MergeStringChannel(ffc1, ffc2, ffc3)


	qf := new(ticker.QuoteFetcher)
	qfc1 := qf.Process(mergedFFChan)
	qfc2 := qf.Process(mergedFFChan)
	qfc3 := qf.Process(mergedFFChan)
	qfc4 := qf.Process(mergedFFChan)
	qfc5 := qf.Process(mergedFFChan)
	mergedQuoteChan := ticker.MergeQuoteChannel(qfc1, qfc2, qfc3, qfc4, qfc5)


	sfp := new(ticker.SecondaryFilterProcessor)
	wellSFChan1 := sfp.Process(mergedQuoteChan)
	wellSFChan2 := sfp.Process(mergedQuoteChan)
	wellSFChan3 := sfp.Process(mergedQuoteChan)

	mergedSFChan := ticker.MergeSFOutputMessageChannel(wellSFChan1, wellSFChan2, wellSFChan3)


	var wellKnown, latestEarnings []string

	for n := range mergedSFChan {
		if (n.WellKnown) {
			wellKnown = append(wellKnown, n.Symbol)
		} else {
			latestEarnings = append(latestEarnings, n.Symbol)
		}
	}

	fmt.Println("Well Known Check")
	for _, v := range wellKnown {
		fmt.Println(v)
	}

	fmt.Println("Latest Earnings Check")
	for _, v := range latestEarnings {
		fmt.Println(v)
	}


	elapsed := time.Since(start)
	log.Printf("Processing took %s",  elapsed)
	log.Printf("Fundamental Filters: Good=%s Bad=%s", ffp.GoodCount, ffp.BadCount)
	log.Printf("Secondary Filters: WK=%s LE=%s", sfp.WellKnowCount, sfp.LatestEarningsCount)
}
