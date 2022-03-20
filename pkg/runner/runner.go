package runner

import (
	"errors"
	"fmt"
	"github.com/james0248/TestDrive.git/pkg/cache"
	"github.com/james0248/TestDrive.git/pkg/request"
	"io"
	"os"
	"os/exec"
)

type runner struct {
	codePath, language string
	website, problemId string
	binaryPath         string
	testCases          []request.TestCase
	logger             *ResultLogger
}

func NewRunner(codePath, language, website, problemId string) *runner {
	r := &runner{
		codePath:  codePath,
		language:  language,
		website:   website,
		problemId: problemId,
	}
	return r
}

func (r *runner) Run(options []string) {
	output, err := r.compile(options)
	if err != nil {
		fmt.Println("Error occured while compiling code")
		fmt.Printf("%s\n", output)
		panic(err)
	}

	err = r.getTestCases()
	if err != nil {
		panic(err)
	}

	r.logger.Start()
	// Run all tests
	for index := range r.testCases {
		go func(i int, tc request.TestCase) {
			err := r.test(i, tc)
			if err != nil {
				fmt.Println(err)
			}
		}(index, r.testCases[index])
	}
	r.logger.WaitForCompletion()

	err = r.removeBinary()
	if err != nil {
		fmt.Println("Error occured while cleaning up")
	}
}

func (r *runner) test(index int, testCase request.TestCase) error {
	cmd := exec.Command(r.binaryPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		// Error handling will be tedious because it occured before execution
		panic(err)
	}

	defer stdin.Close()
	io.WriteString(stdin, testCase.Input)

	output, err := cmd.CombinedOutput()
	if err != nil {
		r.logger.ReportResult(NewRunResult(index, false, "Undiscovered", output, testCase.Output))
		return err
	}

	if string(output) == testCase.Output {
		r.logger.ReportResult(NewRunResult(index, true, "", output, testCase.Output))
	} else {
		r.logger.ReportResult(NewRunResult(index, false, "WA", output, testCase.Output))
	}
	return nil
}

func (r *runner) getTestCases() error {
	// Read test cases from cache
	testCases, err := cache.ReadCache(r.website, r.problemId)
	if err != nil {
		fmt.Println("Error occured while reading cached test cases")
		return err
	}
	// Cache miss
	if testCases == nil {
		testCases, err = request.ParseTestCases(r.website, r.problemId)
		if err != nil {
			fmt.Println("Error occured while fetching test cases from " + r.website)
			return err
		}

		err = cache.WriteCache(r.website, r.problemId, testCases)
		if err != nil {
			fmt.Println("Error occured while writing cache")
			return err
		}
	}
	r.testCases = testCases
	r.logger = NewLogger(len(testCases))
	return nil
}

func (r *runner) compile(options []string) ([]byte, error) {
	compiler := ""
	var compileOptions []string
	if r.language == "C++" {
		compiler = "g++"
		r.binaryPath = "./Main"
		compileOptions = append([]string{r.codePath, "-o", r.binaryPath}, options...)
	}
	// Compile code if neccesary
	cmd := exec.Command(compiler, compileOptions...)
	output, err := cmd.CombinedOutput()
	return output, err
}

func (r *runner) removeBinary() error {
	if _, err := os.Stat(r.binaryPath); err == nil {
		err := os.Remove(r.binaryPath)
		if err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return nil
	} else {
		return err
	}
	return nil
}
