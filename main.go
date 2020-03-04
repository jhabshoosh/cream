package main

import (
	"fmt"
	"sort"
	"github.com/jhabshoo/cream/pipeline/ratios"
	"github.com/jhabshoo/cream/pipeline/ranking"
	"github.com/jhabshoo/cream/pipeline/info"
	"github.com/jhabshoo/cream/pipeline/company"
	"github.com/jhabshoo/cream/pipeline/quote"
	"github.com/jhabshoo/cream/utils"
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


	quoteMap := make(map [string]fmp.CompanyQuote)
	quoteStage := new(quote.QuoteStage)
	quoteStageChan1 := quoteStage.Process(mergedTickerChan, quoteMap)
	quoteStageChan2 := quoteStage.Process(mergedTickerChan, quoteMap)
	quoteStageChan3 := quoteStage.Process(mergedTickerChan, quoteMap)
	quoteStageChan4 := quoteStage.Process(mergedTickerChan, quoteMap)
	quoteStageChan5 := quoteStage.Process(mergedTickerChan, quoteMap)
	mergedQuoteChan := utils.MergeQuoteChannel(quoteStageChan1, quoteStageChan2, quoteStageChan3, quoteStageChan4, quoteStageChan5)


	ratiosStage := new(ratios.RatiosStage)
	ratiosStageChan1 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan2 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan3 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan4 := ratiosStage.Process(mergedQuoteChan)
	ratiosStageChan5 := ratiosStage.Process(mergedQuoteChan)
	mergedRatiosStageChan := utils.MergeStringChannel(ratiosStageChan1, ratiosStageChan2, ratiosStageChan3, ratiosStageChan4, ratiosStageChan5)


	scoreMap := make(map [string]float64)
	rankingStageChan1 := ranking.Rank(mergedRatiosStageChan, quoteMap, scoreMap)
	rankingStageChan2 := ranking.Rank(mergedRatiosStageChan, quoteMap, scoreMap)
	rankingStageChan3 := ranking.Rank(mergedRatiosStageChan, quoteMap, scoreMap)
	rankingStageChan4 := ranking.Rank(mergedRatiosStageChan, quoteMap, scoreMap)
	rankingStageChan5 := ranking.Rank(mergedRatiosStageChan, quoteMap, scoreMap)
	mergedRankingStageChan := utils.MergeRankingScoreChannel(rankingStageChan1, rankingStageChan2, rankingStageChan3, rankingStageChan4, rankingStageChan5)


	var scores []ranking.RankingScore
	for n := range mergedRankingStageChan {
		scores = append(scores, *n)
	}
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	scores = scores[1:100]
	scoreChan := utils.GenerateScoreChannel(scores)


	companyProfileStage := new(company.CompanyProfileStage)
	companyProfileChan := companyProfileStage.Process(scoreChan)
	// mergedCompanyProfileChan := utils.MergeCompanyProfileResponseChannel(companyProfileChan1, companyProfileChan2, companyProfileChan3)

	var profiles []fmp.CompanyProfileResponse
	for n := range companyProfileChan {
		profiles = append(profiles, *n)
	}
	sort.SliceStable(profiles, func(i, j int) bool {
		return scoreMap[profiles[i].Symbol] > scoreMap[profiles[j].Symbol]
	})

	var profilesDeduped []fmp.CompanyProfileResponse
	dedupeMap := make(map[string]int)
	for _, v := range profiles {
		_, ok := dedupeMap[v.Symbol]
		if (!ok) {
			profilesDeduped = append(profilesDeduped, v)
			dedupeMap[v.Symbol] = 1
		}
	}

	fmt.Println("========== REPORT - TOP 50 ==========")
	for _, v := range profilesDeduped[1:51] {
		fmt.Println(fmt.Sprintf("%s | %f", company.ProfileString(&v), scoreMap[v.Symbol]))
	}
	fmt.Println("=====================================")

	elapsed := time.Since(start)
	log.Printf("Processing took %s",  elapsed)
}
