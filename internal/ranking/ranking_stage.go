package ranking

import (
	"log"
	"strconv"
	"sync"

	"github.com/jhabshoosh/cream/internal/base"
	"github.com/jhabshoosh/cream/internal/quote"
	fmp "github.com/jhabshoosh/fmp/pkg/client"
)

type RankingScore struct {
	Symbol string
	Score  float64
}

func (rs RankingScore) GetKey() string {
	return rs.Symbol
}

func (rs RankingScore) SortVal() float64 {
	return rs.Score
}

type scoreMessage struct {
	symbol     string
	cashFlow   fmp.CashFlowStatement
	financials fmp.FinancialStatment
}

func newScoreMessage(symbol string, cashFlow fmp.CashFlowStatement, financials fmp.FinancialStatment) scoreMessage {
	sm := new(scoreMessage)
	sm.symbol = symbol
	sm.cashFlow = cashFlow
	sm.financials = financials
	return *sm
}

func (sm scoreMessage) GetKey() string {
	return sm.symbol
}

func (sm scoreMessage) SortVal() float64 {
	return 0.0
}

type RankingProcessor struct {
	Count    int
	quoteMap *quote.QuoteMap
	scoreMap *ScoreMap
}

type ScoreMap struct {
	mutex sync.Mutex
	Data  map[string]float64
}

func NewScoreMap() *ScoreMap {
	sm := new(ScoreMap)
	sm.Data = make(map[string]float64)
	return sm
}

func NewRankingProcessor(quoteMap *quote.QuoteMap, scoreMap *ScoreMap) *RankingProcessor {
	rp := new(RankingProcessor)
	rp.scoreMap = scoreMap
	rp.quoteMap = quoteMap
	return rp
}

func (rp *RankingProcessor) Filter(m base.Message) bool {
	return true
}

func (rp *RankingProcessor) OutputMessage(im, data base.Message) base.Message {
	d := data.(scoreMessage)
	score := calculateScore(d.financials, d.cashFlow, rp.quoteMap.Data[im.GetKey()])
	return newRankingScore(im.GetKey(), score)
}

func (rp *RankingProcessor) GetData(m base.Message) base.Message {
	cashFlow := getCashFlow(m.GetKey())
	financials := getFinancials(m.GetKey())
	return newScoreMessage(m.GetKey(), cashFlow, financials)
}

func (rp *RankingProcessor) Passed(im, om base.Message) {
	rp.Count++
	rs := om.(RankingScore)
	rp.scoreMap.mutex.Lock()
	defer rp.scoreMap.mutex.Unlock()
	rp.scoreMap.Data[om.GetKey()] = rs.Score
}

func (rp *RankingProcessor) Failed(im, om base.Message) {}

func (ip *RankingProcessor) LogMessage(m base.Message) {
	log.Println("RankingProcessor Received: ", m.GetKey())
}

func newRankingScore(symbol string, score float64) RankingScore {
	rs := new(RankingScore)
	rs.Symbol = symbol
	rs.Score = score
	return *rs
}

func calculateScore(financials fmp.FinancialStatment, cashFlow fmp.CashFlowStatement, q quote.Quote) float64 {
	// So (FCF-Int. Exp.)/Market Cap
	fcf, err := strconv.ParseFloat(cashFlow.FreeCashFlow, 64)
	intExp, err := strconv.ParseFloat(financials.InterestExpense, 64)
	if err != nil {
		return -1
	}
	return (fcf - intExp) / q.MarketCap
}

func getCashFlow(symbol string) fmp.CashFlowStatement {
	fsr, err := fmp.FetchCashFlowStatements(symbol)
	if err != nil {
		log.Println(err)
	}
	if len(fsr.Financials) > 0 {
		return fsr.Financials[0]
	}
	return *new(fmp.CashFlowStatement)
}

func getFinancials(symbol string) fmp.FinancialStatment {
	fsr, err := fmp.FetchFinancialStatements(symbol)
	if err != nil {
		log.Println(err)
	}
	if len(fsr.Financials) > 0 {
		return fsr.Financials[0]
	}
	return *new(fmp.FinancialStatment)
}
