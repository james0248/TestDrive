package main

import (
	"flag"
	"github.com/james0248/TestDrive.git/pkg/cache"
	"github.com/james0248/TestDrive.git/pkg/request"
	"github.com/james0248/TestDrive.git/pkg/runner"
)

func main() {
	var webSite string
	var problemId string
	var codePath string
	flag.StringVar(&webSite, "oj", "", "type of online judge to test")
	flag.StringVar(&problemId, "prob", "", "id of problem to test")
	flag.StringVar(&codePath, "codepath", "", "path to source code")

	flag.Parse()
	if webSite == "" {
		panic("No online judge given")
	}
	if problemId == "" {
		panic("No problem number given")
	}
	if problemId == "" {
		panic("No source code given")
	}
	testCases, err := cache.ReadCache(webSite, problemId)
	if err != nil {
		panic(err)
	}
	testCases, err = request.ParseTestCases(webSite, problemId)
	cache.WriteCache(webSite, problemId, testCases)

	if err != nil {
		panic(err)
	}

	// TODO: Automatically find out language type from extension type
	r := runner.NewRunner(codePath, "C++", webSite, problemId)
	r.Run([]string{"-O2", "-Wall"})
}
