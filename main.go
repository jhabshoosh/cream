package main

import (
	"fmt"
	"github.com/jhabshoo/fmp"
	"time"
	"log"
	s "strings"
	"github.com/jhabshoo/cream/cream"
)

var tickersToProcess = make(chan string, 1000)
var processCount int = 0

func main() {
	start := time.Now()

	companies := fmp.GetSymbolsList()
	var tickers []string
	for _,v := range companies {
		if (!s.Contains(v.Symbol, ".")) {
			tickers = append(tickers, v.Symbol)
		}
	}
	
	infoFetcher := new(cream.InfoFetcher)
	tickerChan := cream.GenerateStringChannel(tickers)
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
	mergedTickerChan := cream.MergeInfoChannel(ic1, ic2, ic3, ic4, ic5, ic6, ic7, ic8, ic9, ic10)


	ffp := new(cream.FundamentalFilterProcessor)
	ffc1 := ffp.Process(mergedTickerChan)
	ffc2 := ffp.Process(mergedTickerChan)
	ffc3 := ffp.Process(mergedTickerChan)
	mergedFFChan := cream.MergeStringChannel(ffc1, ffc2, ffc3)


	qf := new(cream.QuoteFetcher)
	qfc1 := qf.Process(mergedFFChan)
	qfc2 := qf.Process(mergedFFChan)
	qfc3 := qf.Process(mergedFFChan)
	qfc4 := qf.Process(mergedFFChan)
	qfc5 := qf.Process(mergedFFChan)
	mergedQuoteChan := cream.MergeQuoteChannel(qfc1, qfc2, qfc3, qfc4, qfc5)


	sfp := new(cream.SecondaryFilterProcessor)
	wellSFChan1 := sfp.Process(mergedQuoteChan)
	wellSFChan2 := sfp.Process(mergedQuoteChan)
	wellSFChan3 := sfp.Process(mergedQuoteChan)

	mergedSFChan := cream.MergeSFOutputMessageChannel(wellSFChan1, wellSFChan2, wellSFChan3)


	var wellKnown, latestEarnings []string
	for n := range mergedSFChan {
		if (n.WellKnown) {
			wellKnown = append(wellKnown, n.Symbol)
		} else {
			latestEarnings = append(latestEarnings, n.Symbol)
		}
	}

	wellKnownCPFetcher := new(cream.CompanyProfileFetcher)
	latestEarningsCPFetcher := new(cream.CompanyProfileFetcher)
	wellKnownChan := cream.GenerateStringChannel(wellKnown)
	latestEarningsChan := cream.GenerateStringChannel(latestEarnings)

	wellKnownCP1 := wellKnownCPFetcher.Process(wellKnownChan)
	wellKnownCP2 := wellKnownCPFetcher.Process(wellKnownChan)
	wellKnownCP3 := wellKnownCPFetcher.Process(wellKnownChan)
	mergedWellKnownCP := cream.MergeCompanyProfileResponseChannel(wellKnownCP1, wellKnownCP2, wellKnownCP3)

	latestEarningsCP1 := latestEarningsCPFetcher.Process(latestEarningsChan)
	latestEarningsCP2 := latestEarningsCPFetcher.Process(latestEarningsChan)
	latestEarningsCP3 := latestEarningsCPFetcher.Process(latestEarningsChan)
	mergedLatestEarningsCP := cream.MergeCompanyProfileResponseChannel(latestEarningsCP1, latestEarningsCP2, latestEarningsCP3)


	var wellKnownProfiles, latestEarningsProfiles []fmp.CompanyProfileResponse
	for n := range mergedWellKnownCP {
		wellKnownProfiles = append(wellKnownProfiles, *n)
	} 
	for n:= range mergedLatestEarningsCP {
		latestEarningsProfiles = append(latestEarningsProfiles, *n)
	}

	fmt.Println("Well Known Check")
	for _, v := range wellKnownProfiles {
		fmt.Println(v)
	}

	fmt.Println("Latest Earnings Check")
	for _, v := range latestEarningsProfiles {
		fmt.Println(v)
	}


	elapsed := time.Since(start)
	log.Printf("Processing took %s",  elapsed)
	log.Printf("Fundamental Filters: Good=%d Bad=%d", ffp.GoodCount, ffp.BadCount)
	log.Printf("Secondary Filters: WK=%d LE=%d", sfp.WellKnowCount, sfp.LatestEarningsCount)
}
