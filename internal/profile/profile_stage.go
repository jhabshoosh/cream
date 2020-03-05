package profile

import (
	"fmt"
	"log"

	"github.com/jhabshoosh/cream/internal/base"
	"github.com/jhabshoosh/cream/internal/ranking"
	fmp "github.com/jhabshoosh/fmp/pkg/client"
)

type Profile fmp.CompanyProfileResponse

func (p Profile) GetKey() string {
	return p.Symbol
}

func (p Profile) SortVal() float64 {
	return 0
}

type ProfileProcessor struct {
	Count int
}

func (pp *ProfileProcessor) OutputMessage(im, data base.Message) base.Message {
	return pp.GetData(im)
}

func (pp *ProfileProcessor) GetData(m base.Message) base.Message {
	if m == nil {
		return nil
	}
	score := m.(ranking.RankingScore)
	return getCompanyProfile(score.Symbol)
}

func (pp *ProfileProcessor) Filter(m base.Message) bool {
	return true
}

func (pp *ProfileProcessor) Passed(im, om base.Message) {
	pp.Count++
}

func (pp *ProfileProcessor) Failed(im, om base.Message) {}

func (ip *ProfileProcessor) LogMessage(m base.Message) {
	if m != nil {
		log.Println("ProfileProcessor Received:", m.GetKey())
	}
}

func getCompanyProfile(symbol string) Profile {
	cpr, err := fmp.FetchCompanyProfile(symbol)
	if err != nil {
		log.Println(err)
	}
	return Profile(cpr)
}

func ProfileString(p Profile) string {
	return fmt.Sprintf("%s\t%s\t%f\t%s", p.Symbol, p.Profile.CompanyName, p.Profile.Price, p.Profile.Industry)
}
