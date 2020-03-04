package quote

import (
	"sync"
	"log"
	"github.com/jhabshoo/cream/pipeline"
	"github.com/jhabshoo/cream/pipeline/info"
	fmp "github.com/jhabshoo/fmp/client"
)

// QuoteStage fetches from FMP for symbols
type QuoteStage struct {
	GoodCount int
	BadCount int
}

type Quote fmp.CompanyQuote 

func (q Quote) GetKey() string {
	return q.Symbol
}

func (q Quote) SortVal() float64 {
	return 0
}

type QuoteProcessor struct {
	GoodCount int
	BadCount int
	quoteMap *QuoteMap
}

type QuoteMap struct {
	mutex sync.Mutex
	Data map[string]Quote
}

func NewQuoteMap() *QuoteMap {
	qm := new(QuoteMap)
	qm.Data = make(map [string]Quote)
	return qm
}

func NewQuoteProcessor(qm *QuoteMap) *QuoteProcessor {
	qp := new(QuoteProcessor)
	qp.quoteMap = qm
	return qp
}

func (qp *QuoteProcessor) Filter(m pipeline.Message) bool {
	quote := m.(Quote)
	return highLowerFilterRule(quote) && marketCapFilterRule(quote)
}

func (qp *QuoteProcessor) GetData(m pipeline.Message) pipeline.Message {
	i := m.(info.Info)
	quote, _ := getQuote(i.GetKey())
	return quote
}

func (qp *QuoteProcessor) OutputMessage(im, data pipeline.Message) pipeline.Message {
	return qp.GetData(im)
}

func (qp *QuoteProcessor) Passed(im, om pipeline.Message) {
	qp.addQuoteToMap(om.(Quote))
	qp.GoodCount++
}

func (qp *QuoteProcessor) Failed(im, om pipeline.Message) {
	qp.addQuoteToMap(om.(Quote))
	qp.BadCount++
}

func (qp *QuoteProcessor) addQuoteToMap(q Quote) {
	qp.quoteMap.mutex.Lock()
	defer qp.quoteMap.mutex.Unlock()
	qp.quoteMap.Data[q.Symbol] = q
}

func (ip *QuoteProcessor) LogMessage(m pipeline.Message) {
	log.Println("QuoteProcessor Received:", m.GetKey())
}

func getQuote(symbol string) (Quote, error) {
	symbolInput := []string {symbol}
	quoteResponse, err := fmp.FetchCompanyQuote(symbolInput)
	var q Quote
	if (err != nil) {
		return q, err
	}
	if (quoteResponse != nil && len(quoteResponse) > 0) {
		q = Quote(quoteResponse[0])
		return q, nil
	}
	return q, err
}
