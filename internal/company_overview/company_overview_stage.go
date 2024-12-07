package company_overview

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
type CompanyOverviewData struct {
	Symbol       string
	EvOverEbitda float64
	PERatio      float64
	YearHigh     float64
	YearLow      float64
	MarketCap    float64
	Roe          float64
	Roa          float64
}

func (i CompanyOverviewData) GetKey() string {
	return i.Symbol
}

func (i CompanyOverviewData) SortVal() float64 {
	return 0
}

func (i CompanyOverviewData) String() string {
	return fmt.Sprintf("%s | %f", i.Symbol, i.EvOverEbitda)
}

type CompanyOverviewProcessor struct {
	GoodCount int
	BadCount  int
}

func (ip *CompanyOverviewProcessor) Filter(m base.Message) bool {
	cod := m.(CompanyOverviewData)
	return companyOverviewFilter(cod)
}

func (ip *CompanyOverviewProcessor) GetData(m base.Message) base.Message {
	sm := m.(base.StringMessage)
	return getCompanyOverview(sm.GetKey())
}

func (ip *CompanyOverviewProcessor) OutputMessage(im, data base.Message) base.Message {
	return ip.GetData(im)
}

func (ip *CompanyOverviewProcessor) Passed(im, om base.Message) {
	ip.GoodCount++
}

func (ip *CompanyOverviewProcessor) Failed(im, om base.Message) {
	ip.BadCount++
}

func (ip *CompanyOverviewProcessor) LogMessage(m base.Message) {
	log.Println("CompanyOverviewProcessor Received:", m.GetKey())
}

func getCompanyOverview(symbol string) CompanyOverviewData {
	cor, err := av.FetchCompanyOverview(symbol)
	if err != nil {
		log.Println("fetchCompanyOverview:", err)
	}
	cod := new(CompanyOverviewData)
	cod.Symbol = symbol
	cod.EvOverEbitda, err = strconv.ParseFloat(cor.EVToEBITDA, 64)
	if err != nil {
		log.Println("EvOverEbitda:", err)
	}
	cod.PERatio, err = strconv.ParseFloat(cor.PERatio, 64)
	if err != nil {
		log.Println("PERatio:", err)
	}
	cod.YearHigh, err = strconv.ParseFloat(cor.WeekHigh52, 64)
	if err != nil {
		log.Println("YearHigh:", err)
	}
	cod.YearLow, err = strconv.ParseFloat(cor.WeekLow52, 64)
	if err != nil {
		log.Println("YearLow:", err)
	}
	cod.MarketCap, err = strconv.ParseFloat(cor.MarketCapitalization, 64)
	if err != nil {
		log.Println("MarketCap:", err)
	}
	cod.Roe, err = strconv.ParseFloat(cor.ReturnOnEquityTTM, 64)
	if err != nil {
		log.Println("Roe:", err)
	}
	cod.Roa, err = strconv.ParseFloat(cor.ReturnOnAssetsTTM, 64)
	if err != nil {
		log.Println("Roa:", err)
	}

	if err != nil {
		// fmt.Println(err)
	}
	return *cod
}
