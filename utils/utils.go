package utils

import (	
	"github.com/jhabshoo/cream/pipeline/ranking"
	fmp "github.com/jhabshoo/fmp/client"
	"github.com/jhabshoo/cream/pipeline/info"
	"sync"
)

// MergeInfoChannel merges a list of Info channels
func MergeInfoChannel(cs ...<-chan *info.Info) <-chan *info.Info {
	var wg sync.WaitGroup
	out := make(chan *info.Info)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *info.Info) {
			for n := range c {
					out <- n
			}
			wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
			go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}
// MergeQuoteChannel merges a list of CompanyQuote channels
func MergeQuoteChannel(cs ...<-chan *fmp.CompanyQuote) <-chan *fmp.CompanyQuote {
	var wg sync.WaitGroup
	out := make(chan *fmp.CompanyQuote)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *fmp.CompanyQuote) {
			for n := range c {
					out <- n
			}
			wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
			go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}

// MergeStringChannel merges a list of string channels
func MergeStringChannel(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
			for n := range c {
					out <- n
			}
			wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
			go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}

// MergeCompanyProfileResponseChannel merges a list of CompanyProfileResponse channels 
func MergeCompanyProfileResponseChannel(cs ...<-chan *fmp.CompanyProfileResponse) <-chan *fmp.CompanyProfileResponse {
	var wg sync.WaitGroup
	out := make(chan *fmp.CompanyProfileResponse)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *fmp.CompanyProfileResponse) {
			for n := range c {
					out <- n
			}
			wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
			go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}

// MergeRankingScoreChannel merges a list of RankingScore channels 
func MergeRankingScoreChannel(cs ...<-chan *ranking.RankingScore) <-chan *ranking.RankingScore {
	var wg sync.WaitGroup
	out := make(chan *ranking.RankingScore)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *ranking.RankingScore) {
			for n := range c {
					out <- n
			}
			wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
			go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}


// GenerateStringChannel emits values of a []strig to a channel
func GenerateStringChannel(values []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, v := range values {
			out <- v
		}
		close(out)
	}()
	return out
}

// GenerateScoreChannel emits values of a []RankingScore to a channel
func GenerateScoreChannel(values []ranking.RankingScore) <-chan *ranking.RankingScore {
	out := make(chan *ranking.RankingScore)
	go func() {
		for _, v := range values {
			out <- &v
		}
		close(out)
	}()
	return out
}
