package fmp

// KeyMetrics
type KeyMetrics struct {
	EV string `json:"Enterprise Value"`
	EvOverEbitda string `json:"Enterprise Value over EBITDA"`
	PERatio string `json:"PE Ratio"`
}

// KeyMetricsResponse 
type KeyMetricsResponse struct {
	Symbol string `json:"symbol"`
	Metrics []KeyMetrics `json:"metrics"`
}

// Stock
type Stock struct {
	Symbol string
	Name string
	Price float64
	Exchange string
}

type AllCompaniesResponse struct {
	Companies []Stock `json:"symbolsList"`
}

type CompanyQuote struct {
	Symbol string
	Price float64
	ChangesPercentage float64
	Change float64
	DayLow float64
	DayHigh float64
	YearHigh float64
	YearLow float64
	MarketCap float64
	PriceAvg50 float64
	PriceAvg200 float64
	Volume float64
	AvgVolume float64
	Exhange string

}

type CompanyQuoteResponse struct {
	Quotes []CompanyQuote
}