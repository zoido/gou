package config_test

import (
	"errors"
	"testing"

	"github.com/nalgeon/be"

	"github.com/zoido/gou/internal/exec/app/config"
)

func TestParseConfig(t *testing.T) {
	type testCase struct {
		args []string
		want config.Config
	}

	run := func(t *testing.T, tc testCase) {
		// When
		got, err := config.ParseConfig(tc.args)

		// Then
		be.Err(t, err, nil)
		be.Equal(t, tc.want, got)
	}

	testCases := map[string]testCase{
		"no arguments": {
			args: []string{},
		},
		"file flag": {
			args: []string{"-file", "input.txt"},
			want: config.Config{File: "input.txt"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) { run(t, tc) })
	}
}

func TestParseConfig_Help(t *testing.T) {
	// Given
	args := []string{"--help"}

	// When
	_, err := config.ParseConfig(args)
	uErr, isUsage := errors.AsType[config.PrintUsageError](err)

	// Then
	be.True(t, isUsage)
	be.True(t, uErr.Usage != "")
}

func TestParseConfig_Error(t *testing.T) {
	// Given
	args := []string{"-unknown"}

	// When
	_, err := config.ParseConfig(args)

	// Then
	be.True(t, err != nil)
}
