package model_test

import (
	"testing"

	"github.com/nalgeon/be"

	"github.com/zoido/gou/internal/model"
)

func TestBuilder_Build(t *testing.T) {
	type testCase struct {
		url     string
		wantURL string
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		b := model.Builder{URL: tc.url}

		// When
		f, err := b.Build()

		// Then
		be.Err(t, err, nil)
		be.Equal(t, f.URL(), tc.wantURL)
	}

	testCases := map[string]testCase{
		"valid": {
			url:     "https://example.com",
			wantURL: "https://example.com",
		},
		"whitespace is trimmed": {
			url:     "  https://example.com \t\n ",
			wantURL: "https://example.com",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func TestBuilder_Build_Error(t *testing.T) {
	type testCase struct {
		url     string
		wantErr string
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		b := model.Builder{URL: tc.url}

		// When
		_, err := b.Build()

		// Then
		be.Err(t, err, tc.wantErr)
	}

	testCases := map[string]testCase{
		"empty not allowed": {
			url:     "",
			wantErr: "cannot be empty",
		},
		"whitespace-only not allowed": {
			url:     "   ",
			wantErr: "cannot be empty",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func TestFinding_Uninitialised_Panics(t *testing.T) {
	// Given
	var f model.Finding
	var panicked bool

	// When
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_ = f.URL()
	}()

	// Then
	be.True(t, panicked)
}
