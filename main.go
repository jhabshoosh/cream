package main

import (
	"github.com/jhabshoo/cream/pipeline"
	"fmt"
	"github.com/jhabshoo/cream/pipeline/ratios"
	"github.com/jhabshoo/cream/pipeline/ranking"
	"github.com/jhabshoo/cream/pipeline/info"
	"github.com/jhabshoo/cream/pipeline/profile"
	"github.com/jhabshoo/cream/pipeline/quote"
	fmp "github.com/jhabshoo/fmp/client"
	"time"
	"log"
	s "strings"
)

func main() {
	start := time.Now()

	companies := fmp.GetSymbolsList()
	var tickerEnvelope pipeline.Envelope
	for _,v := range companies {
		if (!s.Contains(v.Symbol, ".")) {
			tickerEnvelope = append(tickerEnvelope, pipeline.GetStringMessage(v.Symbol))
		}
	}
	
	tickerChan := pipeline.GenerateChannel(tickerEnvelope)

	infoProcessor := new(info.InfoProcessor)
	infoStage1 := pipeline.Run(infoProcessor, tickerChan)
	infoStage2 := pipeline.Run(infoProcessor, tickerChan)
	infoStage3 := pipeline.Run(infoProcessor, tickerChan)
	infoStage4 := pipeline.Run(infoProcessor, tickerChan)
	infoStage5 := pipeline.Run(infoProcessor, tickerChan)
	infoStage6 := pipeline.Run(infoProcessor, tickerChan)
	infoStage7 := pipeline.Run(infoProcessor, tickerChan)
	infoStage8 := pipeline.Run(infoProcessor, tickerChan)
	infoStage9 := pipeline.Run(infoProcessor, tickerChan)
	infoStage10 := pipeline.Run(infoProcessor, tickerChan)
	mergedInfoChan := pipeline.MergeChannels(infoStage1, infoStage2,infoStage3,infoStage4,infoStage5,infoStage6,infoStage7,infoStage8,infoStage9,infoStage10,)


	quoteMap := quote.NewQuoteMap()
	quoteProcessor := quote.NewQuoteProcessor(quoteMap)
	quoteStage1 := pipeline.Run(quoteProcessor, mergedInfoChan)
	quoteStage2 := pipeline.Run(quoteProcessor, mergedInfoChan)
	quoteStage3 := pipeline.Run(quoteProcessor, mergedInfoChan)
	quoteStage4 := pipeline.Run(quoteProcessor, mergedInfoChan)
	quoteStage5 := pipeline.Run(quoteProcessor, mergedInfoChan)
	mergedQuoteChan := pipeline.MergeChannels(quoteStage1, quoteStage2, quoteStage3, quoteStage4, quoteStage5)


	ratiosProcessor := new(ratios.RatiosProcessor)
	ratiosStage1 := pipeline.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage2 := pipeline.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage3 := pipeline.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage4 := pipeline.Run(ratiosProcessor, mergedQuoteChan)
	ratiosStage5 := pipeline.Run(ratiosProcessor, mergedQuoteChan)
	mergedRatiosChan := pipeline.MergeChannels(ratiosStage1, ratiosStage2, ratiosStage3, ratiosStage4, ratiosStage5)

	scoreMap := ranking.NewScoreMap()
	rankingProcessor := ranking.NewRankingProcessor(quoteMap, scoreMap)
	rankingStage1 := pipeline.Run(rankingProcessor, mergedRatiosChan)
	rankingStage2 := pipeline.Run(rankingProcessor, mergedRatiosChan)
	rankingStage3 := pipeline.Run(rankingProcessor, mergedRatiosChan)
	mergedRankingsChan := pipeline.MergeChannels(rankingStage1, rankingStage2, rankingStage3)


	var scores pipeline.Envelope
	for n := range mergedRankingsChan {
		scores = append(scores, n)
	}
	sortedScores := scores.SortByValue()

	sortedScoresTop100 := sortedScores[1:100]
	scoreChan := pipeline.GenerateChannel(pipeline.Envelope(sortedScoresTop100))


	profileProcessor := new(profile.ProfileProcessor)
	profilesStage1 := pipeline.Run(profileProcessor, scoreChan)
	profilesStage2 := pipeline.Run(profileProcessor, scoreChan)
	profilesStage3 := pipeline.Run(profileProcessor, scoreChan)
	mergedProfilesChan := pipeline.MergeChannels(profilesStage1,profilesStage2, profilesStage3)

	var profiles pipeline.Envelope
	for n := range mergedProfilesChan {
		profiles = append(profiles, n)
	}
	profilesSorted := profiles.Sort()

	var profilesDeduped pipeline.Envelope
	dedupeMap := make(map[string]bool)
	for _, v := range profilesSorted {
		_, ok := dedupeMap[v.GetKey()]
		if (!ok) {
			profilesDeduped = append(profilesDeduped, v)
			dedupeMap[v.GetKey()] = true
		}
	}

	fmt.Println("========== REPORT - TOP 50 ==========")
	for _, v := range profilesDeduped[1:51] {
		company := v.(profile.Profile)
		fmt.Println(fmt.Sprintf("%s | %f", profile.ProfileString(company), scoreMap.Data[v.GetKey()]))
	}
	fmt.Println("=====================================")

	elapsed := time.Since(start)
	log.Printf("Processing took %s",  elapsed)
}
