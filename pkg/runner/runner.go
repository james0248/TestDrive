package runner

import (
	"errors"
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/james0248/TestDrive.git/pkg/cache"
	"github.com/james0248/TestDrive.git/pkg/request"
	"io"
	"os"
	"os/exec"
	"strconv"
)

type runner struct {
	w                  *uilive.Writer
	codePath, language string
	website, problemId string
	binaryPath         string
	testCases          []request.TestCase
}

func NewRunner(codePath, language, website, problemId string) *runner {
	writer := uilive.New()
	r := &runner{
		w:         writer,
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

	r.getTestCases()
	if err != nil {
		panic(err)
	}

	// Run all tests
	results := make(chan *runResult)
	for index := range r.testCases {
		go func(i int, tc request.TestCase) {
			err := r.test(i, tc, results)
			if err != nil {
				fmt.Println(err)
			}
		}(index, r.testCases[index])
	}

	// Observe
	for i := 0; i < len(r.testCases); i++ {
		result := <-results
		if result.ok {
			fmt.Println("Test #" + strconv.Itoa(result.index+1) + " passed")
		} else {
			fmt.Println("Test #" + strconv.Itoa(result.index+1) + " failed: " + result.runtimeErr)
			fmt.Println("Input:")
			fmt.Print(r.testCases[result.index].Input)
			fmt.Println("True Output:")
			fmt.Println(r.testCases[result.index].Output)
			fmt.Println("Your Output:")
			fmt.Print(result.output)
		}
	}

	err = r.removeBinary()
	if err != nil {
		fmt.Println("Error occured while cleaning up")
	}
}

func (r *runner) test(index int, testCase request.TestCase, results chan *runResult) error {
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
		results <- NewRunResult(index, false, "Undiscovered", output)
		return err
	}

	if string(output) == testCase.Output {
		results <- NewRunResult(index, true, "", output)
	} else {
		results <- NewRunResult(index, false, "Wrong Answer", output)
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
