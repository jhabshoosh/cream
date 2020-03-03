package cream

import (
	"fmt"
	"github.com/jhabshoo/fmp"
)


// CompanyProfileFetcher fetches CompanyProfile form channel of symbols
type CompanyProfileFetcher struct {
	Count int
}

func getCompanyProfile(symbol string) fmp.CompanyProfileResponse {
	cpr, err := fmp.FetchCompanyProfile(symbol)
	if (err != nil) {
		fmt.Println(err)
	}
	return cpr
}

// Process Consumes from symbol channel and emits CompanyProfileResponses
func (cpf *CompanyProfileFetcher) Process(in <- chan string) <- chan *fmp.CompanyProfileResponse {
	out := make(chan *fmp.CompanyProfileResponse)
	go func() {
		for v := range in {
			cp := getCompanyProfile(v)
			cpf.Count++
			out <- &cp
		}
		close(out)
	}()
	return out
}