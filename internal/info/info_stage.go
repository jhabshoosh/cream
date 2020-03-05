package info

import (
	"github.com/jhabshoo/cream/internal/base"
	"fmt"
	fmp "github.com/jhabshoo/fmp/pkg/client"
	"log"
	"strconv"
)

var _ = log.Println

// Info dto
type Info struct {
	Symbol       string
	EvOverEbitda float64
	PERatio      float64
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
	BadCount  int
}

func (ip *InfoProcessor) Filter(m base.Message) bool {
	i := m.(Info)
	return multipleFilterRule(i) && peRatioFilterRule(i)
}

func (ip *InfoProcessor) GetData(m base.Message) base.Message {
	sm := m.(base.StringMessage)
	return getInfo(sm.GetKey())
}

func (ip *InfoProcessor) OutputMessage(im, data base.Message) base.Message {
	return ip.GetData(im)
}

func (ip *InfoProcessor) Passed(im, om base.Message) {
	ip.GoodCount++
}

func (ip *InfoProcessor) Failed(im, om base.Message) {
	ip.BadCount++
}

func (ip *InfoProcessor) LogMessage(m base.Message) {
	log.Println("InfoProcessor Received:", m.GetKey())
}

func getInfo(symbol string) Info {
	km, err := fmp.FetchKeyMetrics(symbol)
	info := new(Info)
	info.Symbol = symbol
	if km.Metrics != nil && len(km.Metrics) > 0 {
		info.EvOverEbitda, err = strconv.ParseFloat(km.Metrics[0].EvOverEbitda, 64)
		info.PERatio, err = strconv.ParseFloat(km.Metrics[0].PERatio, 64)
	}
	if err != nil {
		// fmt.Println(err)
	}
	return *info
}
