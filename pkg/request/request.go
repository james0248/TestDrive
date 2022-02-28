package request

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

type TestCase struct {
	Input  string
	Output string
}

func ParseTestCases(problem, webSite string) ([]TestCase, error) {
	doc, err := getHTMLPage(problem, webSite)
	if err != nil {
		fmt.Println("Failed to get HTML page")
		return nil, err
	}

	count, err := getDataCount(doc, webSite)
	if err != nil {
		return nil, err
	}

	inputSelector, outputSelector, err := getSelectors(webSite)
	if err != nil {
		return nil, err
	}

	var testCases []TestCase
	for i := 1; i <= count; i++ {
		input := doc.Find(inputSelector + strconv.Itoa(i)).Text()
		output := doc.Find(outputSelector + strconv.Itoa(i)).Text()
		testCase := TestCase{Input: input, Output: output}
		testCases = append(testCases, testCase)
	}
	return testCases, nil
}

func getHTMLPage(problem, webSite string) (*goquery.Document, error) {
	url, err := getWebsiteUrl(webSite)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url+problem, nil)
	req.Header.Set("User-Agent", "TestDrive") // User-Agent must be set to avoid 403 errors in some websites
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Something went wrong while making HTTP request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Something went wrong while reading HTTP response: ", err)
		return nil, err
	}
	return doc, nil
}

func getWebsiteUrl(webSite string) (string, error) {
	switch webSite {
	case "BOJ":
		return "https://www.acmicpc.net/problem/", nil
	default:
		return "", errors.New("unsupported problem solving website")
	}
}

func getDataCount(doc *goquery.Document, webSite string) (int, error) {
	var selector string
	switch webSite {
	case "BOJ":
		selector = ".sampledata"
	default:
		selector = ""
	}
	count := doc.Find(selector).Size() / 2
	if count == 0 {
		return 0, errors.New("no test cases found")
	}
	return count, nil
}

func getSelectors(webSite string) (string, string, error) {
	switch webSite {
	case "BOJ":
		return "#sample-input-", "#sample-output-", nil
	default:
		return "", "", errors.New("unsupported problem solving website")
	}
}
