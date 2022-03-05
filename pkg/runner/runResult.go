package runner

type runResult struct {
	index      int
	ok         bool
	runtimeErr string
	output     string
}

func NewRunResult(index int, ok bool, runtimeErr string, output []byte) *runResult {
	return &runResult{
		index:      index,
		ok:         ok,
		runtimeErr: runtimeErr,
		output:     string(output),
	}
}
