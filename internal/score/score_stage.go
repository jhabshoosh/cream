package ranking

import (
	"log"
	"strconv"
	"sync"

	av "github.com/jhabshoo/cream/internal/av"
	"github.com/jhabshoo/cream/internal/base"
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
	symbol          string
	companyOverview av.CompanyOverviewResponse
	cashFlow        av.CashFlowResponse
	incomeStatement av.IncomeStatementResponse
}

func newScoreMessage(symbol string, companyOverview av.CompanyOverviewResponse, cashFlow av.CashFlowResponse, incomeStatement av.IncomeStatementResponse) scoreMessage {
	sm := new(scoreMessage)
	sm.symbol = symbol
	sm.companyOverview = companyOverview
	sm.cashFlow = cashFlow
	sm.incomeStatement = incomeStatement
	return *sm
}

func (sm scoreMessage) GetKey() string {
	return sm.symbol
}

func (sm scoreMessage) SortVal() float64 {
	return 0.0
}

type RankingProcessor struct {
	Count int
	// quoteMap *quote.QuoteMap
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

// func NewRankingProcessor(quoteMap *quote.QuoteMap, scoreMap *ScoreMap) *RankingProcessor {
func NewRankingProcessor(scoreMap *ScoreMap) *RankingProcessor {
	rp := new(RankingProcessor)
	rp.scoreMap = scoreMap
	// rp.quoteMap = quoteMap
	return rp
}

func (rp *RankingProcessor) Filter(m base.Message) bool {
	return true
}

func (rp *RankingProcessor) OutputMessage(im, data base.Message) base.Message {
	d := data.(scoreMessage)
	score := calculateScore(d.companyOverview, d.cashFlow, d.incomeStatement)
	return newRankingScore(im.GetKey(), score)
}

func (rp *RankingProcessor) GetData(m base.Message) base.Message {
	companyOverview := getCompanyOverview(m.GetKey())
	cashFlow := getCashFlow(m.GetKey())
	incomeStatement := getIncomeStatement(m.GetKey())
	return newScoreMessage(m.GetKey(), companyOverview, cashFlow, incomeStatement)
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

func calculateScore(companyOverview av.CompanyOverviewResponse, cashFlow av.CashFlowResponse, incomeStatement av.IncomeStatementResponse) float64 {
	// So (FCF-Int. Exp.)/Market Cap
	// fcf, err := strconv.ParseFloat(cashFlow.FreeCashFlow, 64)
	// intExp, err := strconv.ParseFloat(financials.InterestExpense, 64)

	operatingCashFlow, err := strconv.ParseFloat(cashFlow.AnnualReports[0].OperatingCashflow, 64)
	if err != nil {
		return -1
	}
	capitalExpenditures, err := strconv.ParseFloat(cashFlow.AnnualReports[0].CapitalExpenditures, 64)
	if err != nil {
		return -1
	}
	interestExpense, err := strconv.ParseFloat(incomeStatement.AnnualReports[0].InterestExpense, 64)
	if err != nil {
		return -1
	}
	marketCap, err := strconv.ParseFloat(companyOverview.MarketCapitalization, 64)
	if err != nil {
		return -1
	}

	freeCashFlow := operatingCashFlow - capitalExpenditures

	return (freeCashFlow - interestExpense) / marketCap
}

func getCompanyOverview(symbol string) av.CompanyOverviewResponse {
	cor, err := av.FetchCompanyOverview(symbol)
	if err != nil {
		log.Println("fetchCompanyOverview:", err)
	}
	return cor
}

func getCashFlow(symbol string) av.CashFlowResponse {
	cf, err := av.FetchCashFlowStatements(symbol)
	if err != nil {
		log.Println(err)
	}
	return cf
}

func getIncomeStatement(symbol string) av.IncomeStatementResponse {
	is, err := av.FetchIncomeStatement(symbol)
	if err != nil {
		log.Println(err)
	}
	return is
}
