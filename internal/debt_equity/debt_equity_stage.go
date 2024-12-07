package debt_equity

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jhabshoo/cream/internal/av"
	"github.com/jhabshoo/cream/internal/base"
	// fmp "github.com/jhabshoo/fmp/pkg/client"
)

var _ = log.Println

// Info dto
type DebtEquityData struct {
	Symbol                 string
	TotalLiabilities       float64
	TotalShareholderEquity float64
}

func (i DebtEquityData) GetKey() string {
	return i.Symbol
}

func (i DebtEquityData) SortVal() float64 {
	return 0
}

func (i DebtEquityData) String() string {
	return fmt.Sprintf("%s | %f | %f", i.Symbol, i.TotalLiabilities, i.TotalShareholderEquity)
}

type DebtEquityProcessor struct {
	GoodCount int
	BadCount  int
}

func (ip *DebtEquityProcessor) Filter(m base.Message) bool {
	cod := m.(DebtEquityData)
	return debtEquityFilter(cod)
}

func (ip *DebtEquityProcessor) GetData(m base.Message) base.Message {
	sm := m.(base.StringMessage)
	return getDebtEquity(sm.GetKey())
}

func (ip *DebtEquityProcessor) OutputMessage(im, data base.Message) base.Message {
	return ip.GetData(im)
}

func (ip *DebtEquityProcessor) Passed(im, om base.Message) {
	ip.GoodCount++
}

func (ip *DebtEquityProcessor) Failed(im, om base.Message) {
	ip.BadCount++
}

func (ip *DebtEquityProcessor) LogMessage(m base.Message) {
	log.Println("DebtEquityProcessor Received:", m.GetKey())
}

func getDebtEquity(symbol string) DebtEquityData {
	bsr, err := av.FetchBalanceSheet(symbol)
	if err != nil {
		log.Println(err)
	}

	drd := new(DebtEquityData)
	drd.Symbol = symbol
	drd.TotalLiabilities, err = strconv.ParseFloat(bsr.AnnualReports[0].TotalLiabilities, 64)
	if err != nil {
		log.Println(err)
	}
	drd.TotalShareholderEquity, err = strconv.ParseFloat(bsr.AnnualReports[0].TotalShareholderEquity, 64)

	if err != nil {
		// fmt.Println(err)
	}
	return *drd
}
