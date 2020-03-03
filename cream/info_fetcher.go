package cream

import (
	"fmt"
	"github.com/jhabshoo/fmp"
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

// InfoFetcher Fetches KeyMetric info from FMP
type InfoFetcher struct {
	Count int
}

// Process processes symbols from channel
func (cip *InfoFetcher) Process(in <- chan string) <- chan *Info {
	out := make(chan *Info)
	go func() {
		for v := range in {
			ci := getInfo(v)
			cip.Count++
			out <- ci
		}
		close(out)
	}()
	return out
}