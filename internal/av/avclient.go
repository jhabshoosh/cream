package av

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allCompaniesURL = "https://raw.githubusercontent.com/rreichel3/US-Stock-Symbols/refs/heads/main/all/all_tickers.txt"
const companyOverviewURL = "https://www.alphavantage.co/query?function=OVERVIEW"
const cashFlowStatementsURL = "https://www.alphavantage.co/query?function=CASH_FLOW"
const incomeStatementURL = "https://www.alphavantage.co/query?function=INCOME_STATEMENT"
const balanceSheetURL = "https://www.alphavantage.co/query?function=BALANCE_SHEET"
const API_KEY = ""

// KeyMetrics
type KeyMetrics struct {
	EV           string `json:"Enterprise Value"`
	EvOverEbitda string `json:"Enterprise Value over EBITDA"`
	PERatio      string `json:"PE Ratio"`
}

type CashFlowAnnualReport struct {
	FiscalDateEnding                                          string `json:"fiscalDateEnding"`
	ReportedCurrency                                          string `json:"reportedCurrency"`
	OperatingCashflow                                         string `json:"operatingCashflow"`
	PaymentsForOperatingActivities                            string `json:"paymentsForOperatingActivities"`
	ProceedsFromOperatingActivities                           string `json:"proceedsFromOperatingActivities"`
	ChangeInOperatingLiabilities                              string `json:"changeInOperatingLiabilities"`
	ChangeInOperatingAssets                                   string `json:"changeInOperatingAssets"`
	DepreciationDepletionAndAmortization                      string `json:"depreciationDepletionAndAmortization"`
	CapitalExpenditures                                       string `json:"capitalExpenditures"`
	ChangeInReceivables                                       string `json:"changeInReceivables"`
	ChangeInInventory                                         string `json:"changeInInventory"`
	ProfitLoss                                                string `json:"profitLoss"`
	CashflowFromInvestment                                    string `json:"cashflowFromInvestment"`
	CashflowFromFinancing                                     string `json:"cashflowFromFinancing"`
	ProceedsFromRepaymentsOfShortTermDebt                     string `json:"proceedsFromRepaymentsOfShortTermDebt"`
	PaymentsForRepurchaseOfCommonStock                        string `json:"paymentsForRepurchaseOfCommonStock"`
	PaymentsForRepurchaseOfEquity                             string `json:"paymentsForRepurchaseOfEquity"`
	PaymentsForRepurchaseOfPreferredStock                     string `json:"paymentsForRepurchaseOfPreferredStock"`
	DividendPayout                                            string `json:"dividendPayout"`
	DividendPayoutCommonStock                                 string `json:"dividendPayoutCommonStock"`
	DividendPayoutPreferredStock                              string `json:"dividendPayoutPreferredStock"`
	ProceedsFromIssuanceOfCommonStock                         string `json:"proceedsFromIssuanceOfCommonStock"`
	ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `json:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
	ProceedsFromIssuanceOfPreferredStock                      string `json:"proceedsFromIssuanceOfPreferredStock"`
	ProceedsFromRepurchaseOfEquity                            string `json:"proceedsFromRepurchaseOfEquity"`
	ProceedsFromSaleOfTreasuryStock                           string `json:"proceedsFromSaleOfTreasuryStock"`
	ChangeInCashAndCashEquivalents                            string `json:"changeInCashAndCashEquivalents"`
	ChangeInExchangeRate                                      string `json:"changeInExchangeRate"`
	NetIncome                                                 string `json:"netIncome"`
}

type CashFlowResponse struct {
	Symbol        string                 `json:"symbol"`
	AnnualReports []CashFlowAnnualReport `json:"annualReports"`
}

// type AllCompaniesResponse struct {
// 	Companies []Stock `json:"symbolsList"`
// }

type CompanyQuote struct {
	Symbol            string
	Price             float64
	ChangesPercentage float64
	Change            float64
	DayLow            float64
	DayHigh           float64
	YearHigh          float64
	YearLow           float64
	MarketCap         float64
	PriceAvg50        float64
	PriceAvg200       float64
	Volume            float64
	AvgVolume         float64
	Exhange           string
}

type CompanyQuoteResponse struct {
	Quotes []CompanyQuote
}

type CompanyOverviewResponse struct {
	Symbol                     string `json:"Symbol"`
	AssetType                  string `json:"AssetType"`
	Name                       string `json:"Name"`
	Description                string `json:"Description"`
	CIK                        string `json:"CIK"`
	Exchange                   string `json:"Exchange"`
	Currency                   string `json:"Currency"`
	Country                    string `json:"Country"`
	Sector                     string `json:"Sector"`
	Industry                   string `json:"Industry"`
	Address                    string `json:"Address"`
	OfficialSite               string `json:"OfficialSite"`
	FiscalYearEnd              string `json:"FiscalYearEnd"`
	LatestQuarter              string `json:"LatestQuarter"`
	MarketCapitalization       string `json:"MarketCapitalization"`
	EBITDA                     string `json:"EBITDA"`
	PERatio                    string `json:"PERatio"`
	PEGRatio                   string `json:"PEGRatio"`
	BookValue                  string `json:"BookValue"`
	DividendPerShare           string `json:"DividendPerShare"`
	DividendYield              string `json:"DividendYield"`
	EPS                        string `json:"EPS"`
	RevenuePerShareTTM         string `json:"RevenuePerShareTTM"`
	ProfitMargin               string `json:"ProfitMargin"`
	OperatingMarginTTM         string `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          string `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          string `json:"ReturnOnEquityTTM"`
	RevenueTTM                 string `json:"RevenueTTM"`
	GrossProfitTTM             string `json:"GrossProfitTTM"`
	DilutedEPSTTM              string `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY string `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         string `json:"AnalystTargetPrice"`
	AnalystRatingStrongBuy     string `json:"AnalystRatingStrongBuy"`
	AnalystRatingBuy           string `json:"AnalystRatingBuy"`
	AnalystRatingHold          string `json:"AnalystRatingHold"`
	AnalystRatingSell          string `json:"AnalystRatingSell"`
	AnalystRatingStrongSell    string `json:"AnalystRatingStrongSell"`
	TrailingPE                 string `json:"TrailingPE"`
	ForwardPE                  string `json:"ForwardPE"`
	PriceToSalesRatioTTM       string `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           string `json:"PriceToBookRatio"`
	EVToRevenue                string `json:"EVToRevenue"`
	EVToEBITDA                 string `json:"EVToEBITDA"`
	Beta                       string `json:"Beta"`
	WeekHigh52                 string `json:"52WeekHigh"`
	WeekLow52                  string `json:"52WeekLow"`
	DayMovingAverage50         string `json:"50DayMovingAverage"`
	DayMovingAverage200        string `json:"200DayMovingAverage"`
	SharesOutstanding          string `json:"SharesOutstanding"`
	DividendDate               string `json:"DividendDate"`
	ExDividendDate             string `json:"ExDividendDate"`
}

type IncomeStatementAnnualReport struct {
	FiscalDateEnding                  string `json:"fiscalDateEnding"`
	ReportedCurrency                  string `json:"reportedCurrency"`
	GrossProfit                       string `json:"grossProfit"`
	TotalRevenue                      string `json:"totalRevenue"`
	CostOfRevenue                     string `json:"costOfRevenue"`
	CostOfGoodsAndServicesSold        string `json:"costofGoodsAndServicesSold"`
	OperatingIncome                   string `json:"operatingIncome"`
	SellingGeneralAndAdministrative   string `json:"sellingGeneralAndAdministrative"`
	ResearchAndDevelopment            string `json:"researchAndDevelopment"`
	OperatingExpenses                 string `json:"operatingExpenses"`
	InvestmentIncomeNet               string `json:"investmentIncomeNet"`
	NetInterestIncome                 string `json:"netInterestIncome"`
	InterestIncome                    string `json:"interestIncome"`
	InterestExpense                   string `json:"interestExpense"`
	NonInterestIncome                 string `json:"nonInterestIncome"`
	OtherNonOperatingIncome           string `json:"otherNonOperatingIncome"`
	Depreciation                      string `json:"depreciation"`
	DepreciationAndAmortization       string `json:"depreciationAndAmortization"`
	IncomeBeforeTax                   string `json:"incomeBeforeTax"`
	IncomeTaxExpense                  string `json:"incomeTaxExpense"`
	InterestAndDebtExpense            string `json:"interestAndDebtExpense"`
	NetIncomeFromContinuingOperations string `json:"netIncomeFromContinuingOperations"`
	ComprehensiveIncomeNetOfTax       string `json:"comprehensiveIncomeNetOfTax"`
	Ebit                              string `json:"ebit"`
	Ebitda                            string `json:"ebitda"`
	NetIncome                         string `json:"netIncome"`
}

type IncomeStatementResponse struct {
	Symbol        string                        `json:"symbol"`
	AnnualReports []IncomeStatementAnnualReport `json:"annualReports"`
}
type BalanceSheetAnnualReport struct {
	FiscalDateEnding                       string `json:"fiscalDateEnding"`
	ReportedCurrency                       string `json:"reportedCurrency"`
	TotalAssets                            string `json:"totalAssets"`
	TotalCurrentAssets                     string `json:"totalCurrentAssets"`
	CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
	CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
	Inventory                              string `json:"inventory"`
	CurrentNetReceivables                  string `json:"currentNetReceivables"`
	TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
	PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
	AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
	IntangibleAssets                       string `json:"intangibleAssets"`
	IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
	Goodwill                               string `json:"goodwill"`
	Investments                            string `json:"investments"`
	LongTermInvestments                    string `json:"longTermInvestments"`
	ShortTermInvestments                   string `json:"shortTermInvestments"`
	OtherCurrentAssets                     string `json:"otherCurrentAssets"`
	OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
	TotalLiabilities                       string `json:"totalLiabilities"`
	TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
	CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
	DeferredRevenue                        string `json:"deferredRevenue"`
	CurrentDebt                            string `json:"currentDebt"`
	ShortTermDebt                          string `json:"shortTermDebt"`
	TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
	CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
	LongTermDebt                           string `json:"longTermDebt"`
	CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
	LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
	ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
	OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
	OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
	TotalShareholderEquity                 string `json:"totalShareholderEquity"`
	TreasuryStock                          string `json:"treasuryStock"`
	RetainedEarnings                       string `json:"retainedEarnings"`
	CommonStock                            string `json:"commonStock"`
	CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
}

type BalanceSheetResponse struct {
	Symbol        string                     `json:"symbol"`
	AnnualReports []BalanceSheetAnnualReport `json:"annualReports"`
}

func GetSymbolsList() []string {
	res, err := http.Get(allCompaniesURL)
	if err != nil {
		panic(err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	responseString := string(body)
	return strings.Split(responseString, "\n")
}

func FetchCompanyOverview(symbol string) (CompanyOverviewResponse, error) {
	res, err := http.Get(companyOverviewURL + "&symbol=" + symbol + "&apikey=" + API_KEY)
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(body))

	var cor CompanyOverviewResponse
	err = json.Unmarshal(body, &cor)
	if err != nil {
		fmt.Println("err unmarshalling:", err)
	}
	return cor, err
}

func FetchCashFlowStatements(symbol string) (CashFlowResponse, error) {
	res, err := http.Get(cashFlowStatementsURL + "&symbol=" + symbol + "&apikey=" + API_KEY)
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var cfr CashFlowResponse
	err = json.Unmarshal(body, &cfr)
	if err != nil {
		fmt.Println("err unmarshalling:", err)
	}
	return cfr, err
}

func FetchIncomeStatement(symbol string) (IncomeStatementResponse, error) {
	res, err := http.Get(incomeStatementURL + "&symbol=" + symbol + "&apikey=" + API_KEY)
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var isr IncomeStatementResponse
	err = json.Unmarshal(body, &isr)
	if err != nil {
		fmt.Println("err unmarshalling:", err)
	}
	return isr, err
}

func FetchBalanceSheet(symbol string) (BalanceSheetResponse, error) {
	res, err := http.Get(balanceSheetURL + "&symbol=" + symbol + "&apikey=" + API_KEY)
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var bsr BalanceSheetResponse
	err = json.Unmarshal(body, &bsr)
	if err != nil {
		fmt.Println("err unmarshalling:", err)
	}
	return bsr, err
}
