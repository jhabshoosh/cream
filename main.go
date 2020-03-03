package main

import (
	"github.com/jhabshoo/cream/pipeline/ratios"
	"github.com/jhabshoo/cream/pipeline/info"
	"github.com/jhabshoo/cream/pipeline/company"
	"github.com/jhabshoo/cream/pipeline/quote"
	"github.com/jhabshoo/cream/utils"
	"fmt"
	fmp "github.com/jhabshoo/fmp/client"
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
	
	tickerChan := utils.GenerateStringChannel(tickers)


	infoStage := new(info.InfoStage)
	infoStageChan1 := infoStage.Process(tickerChan)
	infoStageChan2 := infoStage.Process(tickerChan)
	infoStageChan3 := infoStage.Process(tickerChan)
	infoStageChan4 := infoStage.Process(tickerChan)
	infoStageChan5 := infoStage.Process(tickerChan)
	infoStageChan6 := infoStage.Process(tickerChan)
	infoStageChan7 := infoStage.Process(tickerChan)
	infoStageChan8 := infoStage.Process(tickerChan)
	infoStageChan9 := infoStage.Process(tickerChan)
	infoStageChan10 := infoStage.Process(tickerChan)
	mergedTickerChan := utils.MergeInfoChannel(
		infoStageChan1, infoStageChan2, infoStageChan3, 
		infoStageChan4, infoStageChan5, infoStageChan6, 
		infoStageChan7, infoStageChan8, infoStageChan9, infoStageChan10)


	quoteStage := new(quote.QuoteStage)
	quoteStageChan1 := quoteStage.Process(mergedTickerChan)
	quoteStageChan2 := quoteStage.Process(mergedTickerChan)
	quoteStageChan3 := quoteStage.Process(mergedTickerChan)
	quoteStageChan4 := quoteStage.Process(mergedTickerChan)
	quoteStageChan5 := quoteStage.Process(mergedTickerChan)
	mergedQuoteChan := utils.MergeQuoteChannel(quoteStageChan1, quoteStageChan2, quoteStageChan3, quoteStageChan4, quoteStageChan5)


	ratiosStage := new(ratios.RatiosStage)
	ratiosStageChan1 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan2 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan3 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan4 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan5 := ratiosStage.Process(mergedQuoteChan)
	mergedRatiosStageChan := utils.MergeStringChannel(ratiosStageChan1, ratiosStageChan2, ratiosStageChan3, ratiosStageChan4, ratiosStageChan5)


	companyProfileStage := new(company.CompanyProfileStage)
	companyProfileChan1 := companyProfileStage.Process(mergedRatiosStageChan)
	companyProfileChan2 := companyProfileStage.Process(mergedRatiosStageChan)
	companyProfileChan3 := companyProfileStage.Process(mergedRatiosStageChan)
	mergedCompanyProfileChan := utils.MergeCompanyProfileResponseChannel(companyProfileChan1, companyProfileChan2, companyProfileChan3)

	for n := range mergedCompanyProfileChan {
		fmt.Println(n)
	}

	elapsed := time.Since(start)
	log.Printf("Processing took %s",  elapsed)
	log.Printf("%d symbols passed filters", companyProfileStage.Count)
}
