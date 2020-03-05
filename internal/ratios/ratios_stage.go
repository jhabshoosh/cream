package ratios

import (
	"log"

	"github.com/jhabshoosh/cream/internal/base"
	fmp "github.com/jhabshoosh/fmp/pkg/client"
)

type Ratio struct {
	Data   fmp.FinancialRatio
	Symbol string
}

func (r Ratio) GetKey() string {
	return r.Symbol
}

func (r Ratio) SortVal() float64 {
	return 0
}

type RatiosProcessor struct {
	GoodCount int
	BadCount  int
}

func (rp *RatiosProcessor) Filter(m base.Message) bool {
	ratio := m.(Ratio)
	return roeFilter(ratio) && roaFilter(ratio) && deFilter(ratio)
}

func (rp *RatiosProcessor) OutputMessage(im, data base.Message) base.Message {
	return base.GetStringMessage(im.GetKey())
}

func (rp *RatiosProcessor) GetData(m base.Message) base.Message {
	return getRatio(m.GetKey())
}

func (rp *RatiosProcessor) Passed(im, om base.Message) {
	rp.GoodCount++
}

func (rp *RatiosProcessor) Failed(im, om base.Message) {
	rp.BadCount++
}

func (ip *RatiosProcessor) LogMessage(m base.Message) {
	log.Println("RatiosProcessor Received:", m.GetKey())
}

func NewRatioFromResponse(r fmp.FinancialRatiosResponse) Ratio {
	ratio := new(Ratio)
	ratio.Symbol = r.Symbol
	if len(r.Ratios) > 0 {
		ratio.Data = r.Ratios[0]
	}
	return *ratio
}

type RatiosStage struct {
	GoodCount int
	BadCount  int
}

func getRatio(symbol string) Ratio {
	frResponse, err := fmp.FetchFinancialRatios(symbol)
	if err != nil {
		log.Println(err)
	}
	return NewRatioFromResponse(frResponse)
}
