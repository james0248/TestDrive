package cache

import (
	"encoding/json"
	"fmt"
	"github.com/james0248/TestDrive.git/pkg/request"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Cache struct {
	Count     int
	TestCases []request.TestCase
}

func newCache(testCases []request.TestCase) *Cache {
	return &Cache{Count: len(testCases), TestCases: testCases}
}

// TODO: Change saving format to JSON
func WriteCache(webSite, problem string, testCases []request.TestCase) error {
	cacheDir, err := getCacheDir(webSite)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}

	cache := newCache(testCases)
	doc, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.json", problem)
	err = ioutil.WriteFile(filepath.Join(cacheDir, filename), doc, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadCache(webSite, problem string) ([]request.TestCase, error) {
	return nil, nil
}

func getCacheDir(webSite string) (string, error) {
	if cacheDir, err := os.UserCacheDir(); err != nil {
		fmt.Println("No local cache directory found: ", err)
		return "", err
	} else {
		p := filepath.Join(cacheDir, "TestDrive", webSite)
		return p, nil
	}
}
