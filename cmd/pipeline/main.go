package main

import (
	"fmt"
	fmp "github.com/jhabshoo/fmp/pkg/client"
	"github.com/jhabshoo/cream/internal/base"
	"github.com/jhabshoo/cream/internal/info"
	"github.com/jhabshoo/cream/internal/quote"
	"github.com/jhabshoo/cream/internal/profile"
	"github.com/jhabshoo/cream/internal/ratios"
	"github.com/jhabshoo/cream/internal/ranking"
	"log"
	s "strings"
	"time"
)

func main() {
	start := time.Now()

	companies := fmp.GetSymbolsList()
	var tickerEnvelope base.Envelope
	for _, v := range companies {
		if !s.Contains(v.Symbol, ".") {
			tickerEnvelope = append(tickerEnvelope, base.GetStringMessage(v.Symbol))
		}
	}

	tickerChan := base.GenerateChannel(tickerEnvelope)

	infoProcessor := new(info.InfoProcessor)
	infoStage1 := base.Run(infoProcessor, tickerChan)
	infoStage2 := base.Run(infoProcessor, tickerChan)
	infoStage3 := base.Run(infoProcessor, tickerChan)
	infoStage4 := base.Run(infoProcessor, tickerChan)
	infoStage5 := base.Run(infoProcessor, tickerChan)
	infoStage6 := base.Run(infoProcessor, tickerChan)
	infoStage7 := base.Run(infoProcessor, tickerChan)
	infoStage8 := base.Run(infoProcessor, tickerChan)
	infoStage9 := base.Run(infoProcessor, tickerChan)
	infoStage10 := base.Run(infoProcessor, tickerChan)
	mergedInfoChan := base.MergeChannels(infoStage1, infoStage2, infoStage3, infoStage4, infoStage5, infoStage6, infoStage7, infoStage8, infoStage9, infoStage10)

	quoteMap := quote.NewQuoteMap()
	quoteProcessor := quote.NewQuoteProcessor(quoteMap)
	quoteStage1 := base.Run(quoteProcessor, mergedInfoChan)
	quoteStage2 := base.Run(quoteProcessor, mergedInfoChan)
	quoteStage3 := base.Run(quoteProcessor, mergedInfoChan)
	quoteStage4 := base.Run(quoteProcessor, mergedInfoChan)
	quoteStage5 := base.Run(quoteProcessor, mergedInfoChan)
	mergedQuoteChan := base.MergeChannels(quoteStage1, quoteStage2, quoteStage3, quoteStage4, quoteStage5)

	ratiosProcessor := new(ratios.RatiosProcessor)
	ratiosStage1 := base.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage2 := base.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage3 := base.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage4 := base.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage5 := base.Run(ratiosProcessor, mergedQuoteChan)
	mergedRatiosChan := base.MergeChannels(ratiosStage1, ratiosStage2, ratiosStage3, ratiosStage4, ratiosStage5)

	scoreMap := ranking.NewScoreMap()
	rankingProcessor := ranking.NewRankingProcessor(quoteMap, scoreMap)
	rankingStage1 := base.Run(rankingProcessor, mergedRatiosChan)
	rankingStage2 := base.Run(rankingProcessor, mergedRatiosChan)
	rankingStage3 := base.Run(rankingProcessor, mergedRatiosChan)
	rankingStage4 := base.Run(rankingProcessor, mergedRatiosChan)
	rankingStage5 := base.Run(rankingProcessor, mergedRatiosChan)
	mergedRankingsChan := base.MergeChannels(rankingStage1, rankingStage2, rankingStage3, rankingStage4, rankingStage5)

	var scores base.Envelope
	for n := range mergedRankingsChan {
		scores = append(scores, n)
	}
	sortedScores := scores.SortByValue()

	scoreChan := base.GenerateChannel(base.Envelope(sortedScores[0:100]))

	profileProcessor := new(profile.ProfileProcessor)
	profilesStage := base.Run(profileProcessor, scoreChan)
	mergedProfilesChan := base.MergeChannels(profilesStage)
	
	var profiles base.Envelope
	for n := range mergedProfilesChan {
		seconds := time.Since(start).Seconds()
		totalSoFar := (profileProcessor.Count + infoProcessor.BadCount + quoteProcessor.BadCount + ratiosProcessor.BadCount )
		log.Printf("Processed %d messages.  %g messages per sec\n", totalSoFar, float64(totalSoFar)/seconds)
		profiles = append(profiles, n)
	}

	var profilesDeduped base.Envelope
	dedupeMap := make(map[string]bool)
	for _, v := range profiles {
		_, ok := dedupeMap[v.GetKey()]
		if !ok {
			profilesDeduped = append(profilesDeduped, v)
			dedupeMap[v.GetKey()] = true
		}
	}

	fmt.Printf("========== REPORT - TOP %v ==========\n", len(profilesDeduped))
	fmt.Println("------STATS-------- OPS STATS --------------")
	fmt.Printf("InfoStage Good: %d Bad %d\n", infoProcessor.GoodCount, infoProcessor.BadCount)
	fmt.Printf("QuoteStage Good: %d Bad %d\n", quoteProcessor.GoodCount, quoteProcessor.BadCount)
	fmt.Printf("RatiosStage Good: %d Bad %d\n", ratiosProcessor.GoodCount, ratiosProcessor.BadCount)
	fmt.Printf("RankingStage Count %d\n", rankingProcessor.Count)
	fmt.Printf("ProfileStage Count %d\n", profileProcessor.Count)
	fmt.Println("---------------------------------------")

	for _, v := range profilesDeduped {
		company := v.(profile.Profile)
		fmt.Printf("%s | %f\n", profile.ProfileString(company), scoreMap.Data[v.GetKey()])
	}
	fmt.Println("=====================================")

	elapsed := time.Since(start)
	total := infoProcessor.GoodCount + infoProcessor.BadCount
	log.Printf("Processing %d messages took %s .  %g messages per sec", total, elapsed, float64(total)/elapsed.Seconds())
}
