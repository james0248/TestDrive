package runner

import (
	"fmt"
	"github.com/james0248/TestDrive.git/pkg/cache"
	"github.com/james0248/TestDrive.git/pkg/request"
	"io"
	"os/exec"
)

func Run(codePath, language, website, problemNumber string, options []string) {
	output, err := compile(codePath, language, options)
	if err != nil {
		fmt.Println("Error occured while compiling code")
		fmt.Printf("%s\n", output)
		panic(err)
	}

	// Read test cases from cache
	testCases, err := cache.ReadCache(website, problemNumber)
	if err != nil {
		fmt.Println("Error occured while reading cached test cases")
		panic(err)
	}
	// Cache miss
	if testCases == nil {
		testCases, err = request.ParseTestCases(website, problemNumber)
		if err != nil {
			fmt.Println("Error occured while fetching test cases from " + website)
			panic(err)
		}

		err = cache.WriteCache(website, problemNumber, testCases)
		if err != nil {
			fmt.Println("Error occured while writing cache")
			panic(err)
		}
	}

	//count := len(testCases)
	for _, tc := range testCases {
		err := test(tc, "./main")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func test(testCase request.TestCase, binaryPath string) error {
	cmd := exec.Command(binaryPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	defer stdin.Close()
	io.WriteString(stdin, testCase.Input)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if string(output) == testCase.Output {
		fmt.Println("Test Passed!")
	} else {
		fmt.Println("Wrong Answer")
	}
	return nil
}

func compile(codePath, language string, options []string) ([]byte, error) {
	compiler := ""
	var compileOptions []string
	if language == "C++" {
		compiler = "g++"
		compileOptions = append([]string{codePath, "-o", "main"}, options...)
	}
	// Compile code if neccesary
	cmd := exec.Command(compiler, compileOptions...)
	output, err := cmd.CombinedOutput()
	return output, err
}
