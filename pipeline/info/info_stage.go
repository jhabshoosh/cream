package info

import (
	"log"
	"github.com/jhabshoo/cream/pipeline"
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

func (i Info) GetKey() string {
	return i.Symbol
}

func (i Info) SortVal() float64 {
	return 0
}

func (i Info) String() string {
	return fmt.Sprintf("%s | %f", i.Symbol, i.EvOverEbitda)
}

type InfoProcessor struct {
	GoodCount int
	BadCount int
}

func (ip *InfoProcessor) Filter(m pipeline.Message) bool {
	i := m.(Info)
	return multipleFilterRule(i) && peRatioFilterRule(i)
}

func (ip *InfoProcessor) GetData(m pipeline.Message) pipeline.Message {
	sm := m.(pipeline.StringMessage)
	return getInfo(sm.GetKey())
}

func (ip *InfoProcessor) OutputMessage(im, data pipeline.Message) pipeline.Message {
	return ip.GetData(im)
}

func (ip *InfoProcessor) Passed(im, om pipeline.Message) {
	ip.GoodCount++
}

func (ip *InfoProcessor) Failed(im, om pipeline.Message) {
	ip.BadCount++
}

func (ip *InfoProcessor) LogMessage(m pipeline.Message) {
	log.Println("InfoProcessor Received:", m.GetKey())
}

func getInfo(symbol string) Info {
	km, err := fmp.FetchKeyMetrics(symbol)
	info := new(Info)
	info.Symbol = symbol
	if (km.Metrics != nil && len(km.Metrics) > 0) {
		info.EvOverEbitda, err = strconv.ParseFloat(km.Metrics[0].EvOverEbitda, 64)
		info.PERatio, err = strconv.ParseFloat(km.Metrics[0].PERatio, 64)
	}
	if (err != nil) {
		// fmt.Println(err)
	}
	return *info
}






























// // InfoStage Fetches KeyMetric info from FMP
// type InfoStage struct {
// 	GoodCount int
// 	BadCount int
// }

// // Process processes symbols from channel
// func (cip *InfoStage) Process(in <- chan string) <- chan *Info {
// 	out := make(chan *Info)
// 	go func() {
// 		for v := range in {
// 			ci := getInfo(v)
// 			if (multipleFilterRule(ci) && peRatioFilterRule(ci)) {
// 				cip.GoodCount++
// 				out <- ci
// 			} else {
// 				cip.BadCount++
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }
