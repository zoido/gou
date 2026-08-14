package scan_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/nalgeon/be"

	"github.com/zoido/gou/internal/model"
	"github.com/zoido/gou/internal/scan"
)

func TestStrict(t *testing.T) {
	type testCase struct {
		input    string
		wantURLs []string
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		find := scan.NewStrictFinder()
		r := strings.NewReader(tc.input)

		// When
		result, err := find(t.Context(), r)

		// Then
		be.Err(t, err, nil)
		be.Equal(t, urlStrings(result), tc.wantURLs)
	}

	testCases := map[string]testCase{
		"no URLs": {
			input:    "just plain text with no links",
			wantURLs: []string{},
		},
		"single URL on single line": {
			input:    "visit https://example.com for more info about example.com",
			wantURLs: []string{"https://example.com"},
		},
		"multiple URLs across multiple lines": {
			input: `
				one: https://example.com/first
				some text that mentions example.com
				and https://example.com/second and http://example.com/third
			`,
			wantURLs: []string{
				"https://example.com/first",
				"https://example.com/second",
				"http://example.com/third",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func TestLax(t *testing.T) {
	type testCase struct {
		input    string
		wantURLs []string
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		find := scan.NewLaxFinder()
		r := strings.NewReader(tc.input)

		// When
		result, err := find(t.Context(), r)

		// Then
		be.Err(t, err, nil)
		be.Equal(t, urlStrings(result), tc.wantURLs)
	}

	testCases := map[string]testCase{
		"no URLs": {
			input:    "just plain text with no links",
			wantURLs: []string{},
		},
		"single URL with scheme": {
			input:    "visit https://example.com for more info",
			wantURLs: []string{"https://example.com"},
		},
		"URL without scheme": {
			input:    "see example.com/path/to/page",
			wantURLs: []string{"example.com/path/to/page"},
		},
		"multiple URLs with and without scheme": {
			input: `
				one: https://example.com/first
				also check example.com/second
				and http://example.com/third
			`,
			wantURLs: []string{
				"https://example.com/first",
				"example.com/second",
				"http://example.com/third",
			},
		},
		"multiple URLs on single line": {
			input:    "example.com and ftp://test.org both work",
			wantURLs: []string{"example.com", "ftp://test.org"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func urlStrings(findings []model.Finding) []string {
	out := make([]string, len(findings))
	for i, f := range findings {
		out[i] = f.URL()
	}
	return out
}

func TestFindFn_ReaderError(t *testing.T) {
	// Given
	wantErr := errors.New("error test sentinel")
	r := iotest.ErrReader(wantErr)
	find := scan.NewStrictFinder()

	// When
	findings, err := find(t.Context(), r)

	// Then
	be.True(t, findings == nil)
	be.Err(t, err, wantErr)
}

func TestFindFn_ContextCancelled(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	r := strings.NewReader("https://example.com")
	find := scan.NewStrictFinder()

	// When
	findings, err := find(ctx, r)

	// Then
	be.True(t, len(findings) == 0)
	be.Err(t, err, context.Canceled)
}
