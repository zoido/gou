package config

import (
	"errors"
	"log/slog"

	"github.com/zoido/yag-config"
)

// Config represents the configuration for the gou application.
type Config struct {
	File   string
	Debug  bool
	Strict bool
	Lax    bool
}

// PrintUsageError is an error type used to indicate that the usage of the application
// should be printed.
type PrintUsageError struct {
	// Usage is the usage information to be printed.
	Usage string
}

func (PrintUsageError) Error() string {
	return yag.ErrHelp.Error()
}

// ParseConfig parses the command-line arguments and environment variables and returns
// the configuration.
func ParseConfig(args []string) (*Config, error) {
	cfg := &Config{}

	y := yag.New()
	cfg.register(y)

	err := y.Parse(args)
	if errors.Is(err, yag.ErrHelp) {
		return nil, PrintUsageError{Usage: y.Usage()}
	}
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LogValue implements the [slog.LogValuer] interface.
func (cfg *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("file", cfg.File),
		slog.Bool("debug", cfg.Debug),
	)
}

func (cfg *Config) register(y *yag.Parser) {
	y.String(&cfg.File, "file", "File where to scan for URLs, if not provided STDIN is used.")
	y.Bool(&cfg.Debug, "debug", "Show debugging output.")
	y.Bool(&cfg.Strict, "strict", "Match only proper URLs with schema.")
	y.Bool(&cfg.Lax, "lax", "Match all URL like strings.")
}
