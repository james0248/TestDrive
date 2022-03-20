package runner

import (
	"fmt"
	"github.com/gosuri/uilive"
	"sync"
)

type ResultLogger struct {
	writer     *uilive.Writer
	runResults []*runResult
	print      chan struct{}
	stop       chan struct{}
	mu         *sync.RWMutex
}

func NewLogger(testCount int) *ResultLogger {
	rr := make([]*runResult, testCount)
	for i := range rr {
		rr[i] = NewEmptyRunResult(i)
	}
	return &ResultLogger{
		writer:     uilive.New(),
		runResults: rr,
		print:      make(chan struct{}, 3),
		stop:       make(chan struct{}),
		mu:         &sync.RWMutex{},
	}
}

func (rl *ResultLogger) Start() {
	rl.writer.Start()
	rl.print <- struct{}{}
	go rl.logResults()
}

func (rl *ResultLogger) WaitForCompletion() {
	<-rl.stop
	return
}

func (rl *ResultLogger) ReportResult(result *runResult) {
	rl.mu.Lock()
	rl.runResults[result.index] = result
	rl.mu.Unlock()
	// Signal to print
	rl.print <- struct{}{}
}

func (rl *ResultLogger) logResults() {
	defer close(rl.print)
	for i := 0; i <= len(rl.runResults); i++ {
		<-rl.print
		for _, r := range rl.runResults {
			w := rl.writer.Newline()
			fmt.Fprintln(w, r.String())
		}
		rl.writer.Flush()
	}
	rl.stop <- struct{}{}
	return
}
