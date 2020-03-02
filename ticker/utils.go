package ticker

import (	
	"github.com/jhabshoo/cream/fmp"
	"sync"
	
)

// MergeTickerInfoChannel merges a list of ticker.Info 
func MergeInfoChannel(cs ...<-chan *Info) <-chan *Info {
	var wg sync.WaitGroup
	out := make(chan *Info)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *Info) {
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
// MergeQuoteChannel merges a list of ticker.Info 
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

// MergeStringChannel merges a list of ticker.Info 
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


// MergeStringChannel merges a list of ticker.Info 
func MergeSFOutputMessageChannel(cs ...<-chan *SecondaryFilterOutputMessage) <-chan *SecondaryFilterOutputMessage {
	var wg sync.WaitGroup
	out := make(chan *SecondaryFilterOutputMessage)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *SecondaryFilterOutputMessage) {
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