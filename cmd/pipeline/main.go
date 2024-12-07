package main

import (
	"fmt"
	"log"
	"time"

	av "github.com/jhabshoo/cream/internal/av"
	"github.com/jhabshoo/cream/internal/base"
	company_overview "github.com/jhabshoo/cream/internal/company_overview"
	// fmp "github.com/jhabshoo/fmp/pkg/client"
)

func main() {
	start := time.Now()

	companies := av.GetSymbolsList()
	var tickerEnvelope base.Envelope
	for _, v := range companies {
		// if !s.Contains(v.Symbol, ".") {
		// 	tickerEnvelope = append(tickerEnvelope, base.GetStringMessage(v.Symbol))
		// }
		tickerEnvelope = append(tickerEnvelope, base.GetStringMessage(v))
	}

	tickerChan := base.GenerateChannel(tickerEnvelope)

	companyOverviewProcessor := new(company_overview.CompanyOverviewProcessor)
	companyOverviewStage1 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage2 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage3 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage4 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage5 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage6 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage7 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage8 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage9 := base.Run(companyOverviewProcessor, tickerChan)
	companyOverviewStage10 := base.Run(companyOverviewProcessor, tickerChan)
	mergedCompanyOverviewChan := base.MergeChannels(companyOverviewStage1, companyOverviewStage2, companyOverviewStage3, companyOverviewStage4, companyOverviewStage5, companyOverviewStage6, companyOverviewStage7, companyOverviewStage8, companyOverviewStage9, companyOverviewStage10)

	for n := range mergedCompanyOverviewChan {
		fmt.Println(n)
	}
	// debtEquityProcessor := new(debt_equity.DebtEquityProcessor)
	// debtEquityStage1 := base.Run(debtEquityProcessor, mergedCompanyOverviewChan)
	// debtEquityStage2 := base.Run(debtEquityProcessor, mergedCompanyOverviewChan)
	// debtEquityStage3 := base.Run(debtEquityProcessor, mergedCompanyOverviewChan)
	// debtEquityStage4 := base.Run(debtEquityProcessor, mergedCompanyOverviewChan)
	// debtEquityStage5 := base.Run(debtEquityProcessor, mergedCompanyOverviewChan)
	// mergedDebtEquityChan := base.MergeChannels(debtEquityStage1, debtEquityStage2, debtEquityStage3, debtEquityStage4, debtEquityStage5)

	// scoreMap := score.NewScoreMap()
	// rankingProcessor := score.NewRankingProcessor(scoreMap)
	// rankingStage1 := base.Run(rankingProcessor, mergedDebtEquityChan)
	// rankingStage2 := base.Run(rankingProcessor, mergedDebtEquityChan)
	// rankingStage3 := base.Run(rankingProcessor, mergedDebtEquityChan)
	// rankingStage4 := base.Run(rankingProcessor, mergedDebtEquityChan)
	// rankingStage5 := base.Run(rankingProcessor, mergedDebtEquityChan)
	// mergedRankingsChan := base.MergeChannels(rankingStage1, rankingStage2, rankingStage3, rankingStage4, rankingStage5)

	// var scores base.Envelope
	// for n := range mergedRankingsChan {
	// 	scores = append(scores, n)
	// }
	// sortedScores := scores.SortByValue()

	// // scoreChan := base.GenerateChannel(base.Envelope(sortedScores[0:100]))

	// for s := range sortedScores {
	// 	fmt.Println(sortedScores[s])
	// }

	// profileProcessor := new(profile.ProfileProcessor)
	// profilesStage := base.Run(profileProcessor, scoreChan)
	// mergedProfilesChan := base.MergeChannels(profilesStage)

	// var profiles base.Envelope
	// for n := range mergedProfilesChan {
	// 	seconds := time.Since(start).Seconds()
	// 	totalSoFar := (profileProcessor.Count + companyOverviewProcessor.BadCount + quoteProcessor.BadCount + ratiosProcessor.BadCount)
	// 	log.Printf("Processed %d messages.  %g messages per sec\n", totalSoFar, float64(totalSoFar)/seconds)
	// 	profiles = append(profiles, n)
	// }

	// var profilesDeduped base.Envelope
	// dedupeMap := make(map[string]bool)
	// for _, v := range profiles {
	// 	_, ok := dedupeMap[v.GetKey()]
	// 	if !ok {
	// 		profilesDeduped = append(profilesDeduped, v)
	// 		dedupeMap[v.GetKey()] = true
	// 	}
	// }

	// fmt.Printf("========== REPORT - TOP %v ==========\n", len(profilesDeduped))
	fmt.Println("------STATS-------- OPS STATS --------------")
	fmt.Printf("CompanyOverviewStage Good: %d Bad %d\n", companyOverviewProcessor.GoodCount, companyOverviewProcessor.BadCount)
	// fmt.Printf("debtEquityStage Good: %d Bad %d\n", debtEquityProcessor.GoodCount, debtEquityProcessor.BadCount)
	// fmt.Printf("RankingStage Count %d\n", rankingProcessor.Count)
	fmt.Println("---------------------------------------")

	// for _, v := range profilesDeduped {
	// 	company := v.(profile.Profile)
	// 	fmt.Printf("%s | %f\n", profile.ProfileString(company), scoreMap.Data[v.GetKey()])
	// }
	fmt.Println("=====================================")

	elapsed := time.Since(start)
	total := companyOverviewProcessor.GoodCount + companyOverviewProcessor.BadCount
	log.Printf("Processing %d messages took %s .  %g messages per sec", total, elapsed, float64(total)/elapsed.Seconds())
}
