package info

import (
	"fmt"
	fmp "github.com/jhabshoo/fmp/client"
	"strconv"
)

// Info dto
type Info struct {
	Symbol string
	EvOverEbitda float64
	PERatio float64
}

func (ci Info) String() string {
	return fmt.Sprintf("%s | %f", ci.Symbol, ci.EvOverEbitda)
}

func getInfo(symbol string) *Info {
	km, err := fmp.FetchKeyMetrics(symbol)
	tickerInfo := new(Info)
	tickerInfo.Symbol = symbol
	if (km.Metrics != nil && len(km.Metrics) > 0) {
		tickerInfo.EvOverEbitda, err = strconv.ParseFloat(km.Metrics[0].EvOverEbitda, 64)
		tickerInfo.PERatio, err = strconv.ParseFloat(km.Metrics[0].PERatio, 64)
	}
	if (err != nil) {
		// fmt.Println(err)
	}
	return tickerInfo
}

// InfoStage Fetches KeyMetric info from FMP
type InfoStage struct {
	GoodCount int
	BadCount int
}

// Process processes symbols from channel
func (cip *InfoStage) Process(in <- chan string) <- chan *Info {
	out := make(chan *Info)
	go func() {
		for v := range in {
			ci := getInfo(v)
			if (multipleFilterRule(ci) && peRatioFilterRule(ci)) {
				cip.GoodCount++
				out <- ci
			} else {
				cip.BadCount++
			}
		}
		close(out)
	}()
	return out
}
