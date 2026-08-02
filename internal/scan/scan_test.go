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

func TestFindURLs(t *testing.T) {
	type testCase struct {
		input    string
		wantURLs []string
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		r := strings.NewReader(tc.input)

		// When
		result, err := scan.FindURLs(t.Context(), r)

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
			input:    "visit https://example.com for more info",
			wantURLs: []string{"https://example.com"},
		},
		"multiple URLs across multiple lines": {
			input: `
				one: https://example.com/first
				some text
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

func urlStrings(findings []model.Finding) []string {
	out := make([]string, len(findings))
	for i, f := range findings {
		out[i] = f.URL()
	}
	return out
}

func TestFindURLs_ReaderError(t *testing.T) {
	// Given
	wantErr := errors.New("error test sentinel")
	r := iotest.ErrReader(wantErr)

	// When
	findings, err := scan.FindURLs(t.Context(), r)

	// Then
	be.True(t, findings == nil)
	be.Err(t, err, wantErr)
}

func TestFindURLs_ContextCancelled(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	r := strings.NewReader("https://example.com")

	// When
	findings, err := scan.FindURLs(ctx, r)

	// Then
	be.True(t, len(findings) == 0)
	be.Err(t, err, context.Canceled)
}
