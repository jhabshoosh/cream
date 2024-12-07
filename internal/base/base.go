package base

import (
	"log"
	"sort"
	"sync"
)

type Processor interface {
	Filter(m Message) bool
	OutputMessage(im, data Message) Message
	GetData(m Message) Message
	Passed(im, om Message)
	Failed(im, om Message)
	LogMessage(m Message)
}

type Message interface {
	GetKey() string
	SortVal() float64
}

type Envelope []Message

func (e Envelope) SortByValue() Envelope {
	out := e
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].SortVal() > out[j].SortVal()
	})
	return out
}

func (e Envelope) Sort() Envelope {
	out := e
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].GetKey() > out[j].GetKey()
	})
	return out
}

type StringMessage struct {
	key string
}

func (sm StringMessage) GetKey() string {
	return sm.key
}

func (sm StringMessage) SortVal() float64 {
	return 0
}

func GetStringMessage(str string) StringMessage {
	sm := new(StringMessage)
	sm.key = str
	return *sm
}

func Run(p Processor, in <-chan Message) chan Message {
	out := make(chan Message)
	go func() {
		for im := range in {
			p.LogMessage(im)
			data := p.GetData(im)
			om := p.OutputMessage(im, data)
			if p.Filter(data) {
				p.Passed(im, om)
				if om != nil {
					out <- om
				} else {
					log.Println("Nil message filtered out of pipeline")
				}
			} else {
				p.Failed(im, om)
			}
		}
		close(out)
	}()
	return out
}

func GenerateChannel(e Envelope) <-chan Message {
	out := make(chan Message)
	go func() {
		for _, v := range e {
			out <- v
		}
		close(out)
	}()
	return out
}

func MergeChannels(cs ...<-chan Message) <-chan Message {
	var wg sync.WaitGroup
	out := make(chan Message)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Message) {
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
