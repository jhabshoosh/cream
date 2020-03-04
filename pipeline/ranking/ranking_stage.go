package ranking

import (
	"log"
	"strconv"
	fmp "github.com/jhabshoo/fmp/client"
)

type RankingScore struct {
	Symbol string
	Score float64
}

func newRangingScore(symbol string, score float64) *RankingScore {
	rs := new(RankingScore)
	rs.Symbol = symbol
	rs.Score = score
	return rs
}

func calculateScore(financials fmp.FinancialStatment, cashFlow fmp.CashFlowStatement, quote fmp.CompanyQuote) float64 {
	// So (FCF-Int. Exp.)/Market Cap
	fcf, err := strconv.ParseFloat(cashFlow.FreeCashFlow, 64)
	intExp, err := strconv.ParseFloat(financials.InterestExpense, 64)
	if (err != nil) {
		return -1
	}
	return (fcf - intExp) / quote.MarketCap
}

func getCashFlow(symbol string) fmp.CashFlowStatement {
	fsr, err := fmp.FetchCashFlowStatements(symbol)
	if (err != nil) {
		log.Println(err)
	}
	if (len(fsr.Financials) > 0) {
		return fsr.Financials[0]
	}
	return *new(fmp.CashFlowStatement)
}

func getFinancials(symbol string) fmp.FinancialStatment {
	fsr, err := fmp.FetchFinancialStatements(symbol)
	if (err != nil) {
		log.Println(err)
	}
	if (len(fsr.Financials) > 0) {
		return fsr.Financials[0]
	}
	return *new(fmp.FinancialStatment)
}

func Rank(in <- chan string, quoteMap map[string]fmp.CompanyQuote, scoreMap map[string]float64) <- chan *RankingScore {
	out := make(chan *RankingScore)
	go func() {
		for v := range in {
			cashFlow := getCashFlow(v)
			financials := getFinancials(v)
			q := quoteMap[v]
			score := calculateScore(financials, cashFlow, q)
			scoreMap[v] = score
			rs := newRangingScore(v, score)
			out <- rs
		}
		close(out)
	}()
	return out
}