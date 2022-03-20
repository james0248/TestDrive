package runner

import (
	"fmt"
)

type runResult struct {
	index      int
	running    bool
	ok         bool
	runtimeErr string
	output     string
	expected   string
}

func NewEmptyRunResult(index int) *runResult {
	return &runResult{
		index:   index,
		running: true,
	}
}

func NewRunResult(index int, ok bool, runtimeErr string, output []byte, expected string) *runResult {
	return &runResult{
		index:      index,
		running:    false,
		ok:         ok,
		runtimeErr: runtimeErr,
		output:     string(output),
		expected:   expected,
	}
}

func (r runResult) String() string {
	// Still testing
	if r.running {
		return fmt.Sprintf("Test #%d running...", r.index+1)
	}

	// Done testing, test success
	if r.ok {
		return fmt.Sprintf("Test #%d: Passed!", r.index+1)
	}
	// Test failed
	str := ""
	switch r.runtimeErr {
	case "WA":
		str += fmt.Sprintf("Test #%d: Wrong Answer\n", r.index+1)
		str += fmt.Sprintf("[Expected Output]\n")
		str += fmt.Sprint(r.expected)
		str += fmt.Sprintf("[Your Output]\n")
		str += fmt.Sprint(r.output)
		return str
	case "UD":
		str += fmt.Sprintf("Test #%d: Runtime Error\n", r.index+1)
	}
	return str
}
