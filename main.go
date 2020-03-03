package main

import (
	"github.com/jhabshoo/cream/utils"
	"github.com/jhabshoo/cream/info"
	"github.com/jhabshoo/cream/company"
	"github.com/jhabshoo/cream/quote"
	"github.com/jhabshoo/cream/filter"
	"fmt"
	"github.com/jhabshoo/fmp"
	"time"
	"log"
	s "strings"
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
	
	infoFetcher := new(info.InfoFetcher)
	tickerChan := utils.GenerateStringChannel(tickers)
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
	mergedTickerChan := utils.MergeInfoChannel(ic1, ic2, ic3, ic4, ic5, ic6, ic7, ic8, ic9, ic10)


	ffp := new(filter.FundamentalFilterProcessor)
	ffc1 := ffp.Process(mergedTickerChan)
	ffc2 := ffp.Process(mergedTickerChan)
	ffc3 := ffp.Process(mergedTickerChan)
	mergedFFChan := utils.MergeStringChannel(ffc1, ffc2, ffc3)


	qf := new(quote.QuoteFetcher)
	qfc1 := qf.Process(mergedFFChan)
	qfc2 := qf.Process(mergedFFChan)
	qfc3 := qf.Process(mergedFFChan)
	qfc4 := qf.Process(mergedFFChan)
	qfc5 := qf.Process(mergedFFChan)
	mergedQuoteChan := utils.MergeQuoteChannel(qfc1, qfc2, qfc3, qfc4, qfc5)


	sfp := new(filter.SecondaryFilterProcessor)
	wellSFChan1 := sfp.Process(mergedQuoteChan)
	wellSFChan2 := sfp.Process(mergedQuoteChan)
	wellSFChan3 := sfp.Process(mergedQuoteChan)

	mergedSFChan := utils.MergeSFOutputMessageChannel(wellSFChan1, wellSFChan2, wellSFChan3)


	var wellKnown, latestEarnings []string
	for n := range mergedSFChan {
		if (n.WellKnown) {
			wellKnown = append(wellKnown, n.Symbol)
		} else {
			latestEarnings = append(latestEarnings, n.Symbol)
		}
	}

	wellKnownCPFetcher := new(company.CompanyProfileFetcher)
	latestEarningsCPFetcher := new(company.CompanyProfileFetcher)
	wellKnownChan := utils.GenerateStringChannel(wellKnown)
	latestEarningsChan := utils.GenerateStringChannel(latestEarnings)

	wellKnownCP1 := wellKnownCPFetcher.Process(wellKnownChan)
	wellKnownCP2 := wellKnownCPFetcher.Process(wellKnownChan)
	wellKnownCP3 := wellKnownCPFetcher.Process(wellKnownChan)
	mergedWellKnownCP := utils.MergeCompanyProfileResponseChannel(wellKnownCP1, wellKnownCP2, wellKnownCP3)

	latestEarningsCP1 := latestEarningsCPFetcher.Process(latestEarningsChan)
	latestEarningsCP2 := latestEarningsCPFetcher.Process(latestEarningsChan)
	latestEarningsCP3 := latestEarningsCPFetcher.Process(latestEarningsChan)
	mergedLatestEarningsCP := utils.MergeCompanyProfileResponseChannel(latestEarningsCP1, latestEarningsCP2, latestEarningsCP3)


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
