package company

import (
	"log"
	"fmt"
	"github.com/jhabshoo/cream/pipeline/ranking"
	fmp "github.com/jhabshoo/fmp/client"
)


// CompanyProfileStage fetches CompanyProfile form channel of symbols
type CompanyProfileStage struct {
	Count int
}

func getCompanyProfile(symbol string) fmp.CompanyProfileResponse {
	cpr, err := fmp.FetchCompanyProfile(symbol)
	if (err != nil) {
		log.Println(err)
	}
	return cpr
}

// Process Consumes from symbol channel and emits CompanyProfileResponses
func (cpf *CompanyProfileStage) Process(in <- chan *ranking.RankingScore) <- chan *fmp.CompanyProfileResponse {
	out := make(chan *fmp.CompanyProfileResponse)
	go func() {
		for v := range in {
			cp := getCompanyProfile(v.Symbol)
			cpf.Count++
			out <- &cp
		}
		close(out)
	}()
	return out
}

func ProfileString(p *fmp.CompanyProfileResponse) string {
	return fmt.Sprintf("%s\t%s\t%f\t%s", p.Symbol, p.Profile.CompanyName, p.Profile.Price, p.Profile.Industry)
}