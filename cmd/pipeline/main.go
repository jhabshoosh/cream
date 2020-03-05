package main

import (
	"fmt"
	"log"
	s "strings"
	"time"

	"github.com/jhabshoo/cream/internal/base"
	"github.com/jhabshoo/cream/internal/info"
	"github.com/jhabshoo/cream/internal/profile"
	"github.com/jhabshoo/cream/internal/quote"
	"github.com/jhabshoo/cream/internal/ranking"
	"github.com/jhabshoo/cream/internal/ratios"
	fmp "github.com/jhabshoo/fmp/pkg/client"
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
	infoStage := base.Run(infoProcessor, tickerChan, 10)
	mergedInfoChan := base.MergeChannels(infoStage)

	quoteMap := quote.NewQuoteMap()
	quoteProcessor := quote.NewQuoteProcessor(quoteMap)
	quoteStage := base.Run(quoteProcessor, mergedInfoChan, 10)
	mergedQuoteChan := base.MergeChannels(quoteStage)

	ratiosProcessor := new(ratios.RatiosProcessor)
	ratiosStage := base.Run(ratiosProcessor, mergedQuoteChan, 10)
	mergedRatiosChan := base.MergeChannels(ratiosStage)

	scoreMap := ranking.NewScoreMap()
	rankingProcessor := ranking.NewRankingProcessor(quoteMap, scoreMap)
	rankingStage := base.Run(rankingProcessor, mergedRatiosChan, 5)
	mergedRankingsChan := base.MergeChannels(rankingStage)

	var scores base.Envelope
	for n := range mergedRankingsChan {
		scores = append(scores, n)
	}
	sortedScores := scores.SortByValue()

	scoreChan := base.GenerateChannel(base.Envelope(sortedScores[0:100]))

	profileProcessor := new(profile.ProfileProcessor)
	profilesStage := base.Run(profileProcessor, scoreChan, 1)
	mergedProfilesChan := base.MergeChannels(profilesStage)

	var profiles base.Envelope
	for n := range mergedProfilesChan {
		seconds := time.Since(start).Seconds()
		totalSoFar := (profileProcessor.Count + infoProcessor.BadCount + quoteProcessor.BadCount + ratiosProcessor.BadCount)
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
