package main

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_checkResponseCode(t *testing.T) {
	tests := []struct {
		name string

		url string

		mockedStatusCode int

		wantErr string
	}{
		{
			name:             "OK",
			mockedStatusCode: http.StatusOK,
		},
		{
			name:    "Error during GET",
			url:     "htp://example.com", // invalid protocol instead of mocking status code
			wantErr: `while making http request: Get "htp://example.com": unsupported protocol scheme "htp"`,
		},
		{
			name:             "Invalid status code",
			mockedStatusCode: http.StatusNotFound,
			wantErr:          `invalid http response code: 404`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Mock response
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockedStatusCode)
			}))
			defer ts.Close()
			testURL := ts.URL
			// If url is set for this test case, use instead of the mock server
			if tt.url != "" {
				testURL = tt.url
			}
			err := checkResponseCode(testURL)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name       string
		env        string
		wantStdout []string
		wantErr    string
	}{
		{
			name: "OK",
			env:  "testdata/urls.txt",
			wantStdout: []string{
				`URL: htp://example.com, Err: while making http request: Get "htp://example.com": unsupported protocol scheme "htp"`,
				`URL: http://example.com, OK`,
				`URL: http://httpstat.us/200, OK`,
				`URL: http://httpstat.us/404, Err: invalid http response code: 404`,
				`URL: http://httpstat.us/301, OK`,
			},
		},
		{
			name:    "Error while reading env",
			wantErr: "while loading environment: FILE_PATH environment variable is not set",
		},
		{
			name:    "Error while parsing urls file",
			env:     "testdata/na.txt",
			wantErr: "while opening file: open testdata/na.txt: no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env
			if tt.env != "" {
				require.NoError(t, os.Setenv("FILE_PATH", tt.env))
				defer func() {
					require.NoError(t, os.Unsetenv("FILE_PATH"))
				}()
			}

			// Redirect stdout to file
			stdout, err := os.Create("testdata/stdout.txt")
			require.NoError(t, err)
			defer stdout.Close()

			// Cleanup output file
			defer func() {
				require.NoError(t, os.Remove("testdata/stdout.txt"))
			}()

			os.Stdout = stdout

			// Run
			err = run()

			// Compare stdout and err with expected
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			stdoutLines, err := readLines("testdata/stdout.txt")
			require.NoError(t, err)
			require.ElementsMatch(t, tt.wantStdout, stdoutLines)
		})
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Test_parseURLFile(t *testing.T) {
	tests := []struct {
		name     string
		r        io.Reader
		wantURLs []string
		wantErr  string
	}{
		{
			name: "OK",
			r:    strings.NewReader("http://example.com\nhttp://httpstat.us/200\nhttp://httpstat.us/404\nhttp://httpstat.us/301\n"),
			wantURLs: []string{
				"http://example.com",
				"http://httpstat.us/200",
				"http://httpstat.us/404",
				"http://httpstat.us/301",
			},
		},

		{
			name:     "OK empty",
			r:        strings.NewReader(""),
			wantURLs: []string{},
		},

		{
			name: "OK empty lines",
			r:    strings.NewReader("\n\n"),
			wantURLs: []string{
				"",
				"",
			},
		},

		{
			name:    "io error",
			r:       &errorReader{err: errors.New("io error")},
			wantErr: "while reading file: io error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURLFile(tt.r)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.wantURLs, got)
		})
	}
}

// errorReader is a reader that always returns an error
type errorReader struct {
	err error
}

// Read always returns an error
func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func Test_processURLs(t *testing.T) {
	// process the function return the url via the error for a check
	processFn := func(url string) error {
		return errors.New(url)
	}
	tests := []struct {
		name string
		urls []string
		want []statusCodeCheckResult
	}{
		{
			name: "OK",
			urls: []string{
				"http://example.com",
				"http://httpstat.us/200",
			},
			want: []statusCodeCheckResult{
				{
					url: "http://example.com",
					err: errors.New("http://example.com"),
				},
				{
					url: "http://httpstat.us/200",
					err: errors.New("http://httpstat.us/200"),
				},
			},
		},
		{
			name: "OK no urls",
			urls: []string{},
			want: []statusCodeCheckResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processURLs(tt.urls, processFn)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
