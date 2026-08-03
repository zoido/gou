package dispatch_test

import (
	"errors"
	"testing"

	"github.com/nalgeon/be"

	"github.com/zoido/gou/internal/dispatch"
	"github.com/zoido/gou/internal/model"
)

func TestDispatcher_Dispatch(t *testing.T) {
	type testCase struct {
		dispatcher dispatch.Dispatcher
		key        string
		wantDoQuit bool
		wantErr    any
	}

	run := func(t *testing.T, tc testCase) {
		// Given
		f, err := model.Builder{URL: "https://example.com"}.Build()
		if err != nil {
			t.Fatalf("creating finding: %v", err)
		}

		// When
		gotDoQuit, err := tc.dispatcher.Dispatch(tc.key, f)

		// Then
		be.Err(t, err, tc.wantErr)
		be.Equal(t, gotDoQuit, tc.wantDoQuit)
	}

	testCases := map[string]testCase{
		"known key returns action result": {
			dispatcher: dispatch.Dispatcher{
				"q": func(model.Finding) (bool, error) { return true, nil },
			},
			key:        "q",
			wantDoQuit: true,
		},
		"known key returns false result": {
			dispatcher: dispatch.Dispatcher{
				"o": func(model.Finding) (bool, error) { return false, nil },
			},
			key:        "o",
			wantDoQuit: false,
		},
		"unknown key is a no-op": {
			dispatcher: dispatch.Dispatcher{
				"q": func(model.Finding) (bool, error) { return true, errors.New("an error") },
			},
			key:        "x",
			wantDoQuit: false,
			wantErr:    nil,
		},
		"empty dispatcher is a no-op": {
			dispatcher: dispatch.Dispatcher{},
			key:        "q",
			wantDoQuit: false,
		},
		"action error is propagated": {
			dispatcher: dispatch.Dispatcher{
				"e": func(model.Finding) (bool, error) { return false, errors.New("test error") },
			},
			key:        "e",
			wantDoQuit: false,
			wantErr:    "test error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func TestDispatcher_Dispatch_PassesFinding(t *testing.T) {
	// Given
	f, err := model.Builder{URL: "https://example.com"}.Build()
	if err != nil {
		t.Fatal(err)
	}
	var gotURL string
	d := dispatch.Dispatcher{
		"o": func(f model.Finding) (bool, error) {
			gotURL = f.URL()
			return false, nil
		},
	}

	// When
	_, err = d.Dispatch("o", f)

	// Then
	be.Err(t, err, nil)
	be.Equal(t, gotURL, "https://example.com")
}
