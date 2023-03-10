package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Unexpected error: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	filePath, outputFilePath, err := loadEnv()
	if err != nil {
		return fmt.Errorf("while loading environment: %w", err)
	}

	// open file
	inputF, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("while opening input file: %w", err)
	}
	defer inputF.Close()

	urls, err := parseURLFile(inputF)
	if err != nil {
		return fmt.Errorf("while parsing file: %w", err)
	}

	results := processURLs(urls, checkResponseCode)

	outputF, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("while opening output file: %w", err)
	}
	defer outputF.Close()

	printResults(outputF, results)

	return nil
}

func loadEnv() (string, string, error) {
	// Read file location from environment variable
	input := os.Getenv("INPUT_PATH")
	if input == "" {
		return "", "", fmt.Errorf("INPUT_PATH environment variable is not set")
	}

	// Read file location from environment variable
	output := os.Getenv("OUTPUT_PATH")
	if output == "" {
		return "", "", fmt.Errorf("OUTPUT_PATH environment variable is not set")
	}

	return input, output, nil
}

// parseURLFile reads a file and returns a slice of urls.
// The file read is expected to be a list of real-world http urls is separated by a newline (each url in a single line).
func parseURLFile(r io.Reader) ([]string, error) {
	// Read each line into a slice of urls
	var (
		scanner = bufio.NewScanner(r)
		urls    = make([]string, 0)
	)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("while reading file: %w", err)
	}

	// return slice of urls
	return urls, nil
}

// statusCodeCheckResult is a struct that contains the url and the error (if any) of the http request.
type statusCodeCheckResult struct {
	url string
	err error
}

// processURLs makes a http request to each url and returns a slice of statusCodeCheckResult.
func processURLs(urls []string, processFn func(string) error) []statusCodeCheckResult {
	var (
		resultsCh = make(chan statusCodeCheckResult, len(urls))
		wg        = &sync.WaitGroup{}
	)

	// start a goroutine for each url
	// TODO: limit the number of concurrent requests if needed
	for i := range urls {
		url := urls[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			resultsCh <- statusCodeCheckResult{
				url: url,
				err: processFn(url),
			}
		}()
	}

	// wait for all goroutines to finish and close results channel
	wg.Wait()
	close(resultsCh)

	// read results from channel
	var result []statusCodeCheckResult
	for r := range resultsCh {
		result = append(result, r)
	}

	return result
}

// checkResponseCode makes a http request to the given url and returns an error if the response code is not 200.
// TODO: check redirects and other valid response codes.
func checkResponseCode(url string) error {
	// make http request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("while making http request: %w", err)
	}
	defer resp.Body.Close()

	// check response code is valid
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid http response code: %d", resp.StatusCode)
	}

	return nil
}

// printResults prints the results to the given writer.
// TODO: allow printing in machine-readable format (e.g. add env option to print in json format)
func printResults(w io.Writer, results []statusCodeCheckResult) {
	for _, r := range results {
		if r.err != nil {
			fmt.Fprintf(w, "URL: %s, Err: %s\n", r.url, r.err.Error())
			continue
		}
		fmt.Fprintf(w, "URL: %s, OK\n", r.url)
	}
}
