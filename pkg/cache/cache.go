package cache

import (
	"fmt"
	"github.com/james0248/TestDrive.git/pkg/request"
	"os"
	"path/filepath"
)

func WriteCache(webSite, problem string, testCases []request.TestCase) error {
	cacheDir, err := getCacheDir(webSite, problem)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}

	for i, tc := range testCases {
		infName := fmt.Sprintf("in-%d", i+1)
		inf, err := os.Create(filepath.Join(cacheDir, infName))
		if err != nil {
			panic(err)
		}

		_, err = inf.WriteString(tc.Input)
		if err != nil {
			panic(err)
		}

		outfName := fmt.Sprintf("out-%d", i+1)
		outf, err := os.Create(filepath.Join(cacheDir, outfName))
		if err != nil {
			panic(err)
		}

		_, err = outf.WriteString(tc.Output)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func ReadCache(webSite, problem string) ([]request.TestCase, error) {
	return nil, nil
}

func getCacheDir(webSite, problem string) (string, error) {
	if cacheDir, err := os.UserCacheDir(); err != nil {
		fmt.Println("No local cache directory found: ", err)
		return "", err
	} else {
		p := filepath.Join(cacheDir, "TestDrive", webSite, problem)
		return p, nil
	}
}
