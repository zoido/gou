package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/zoido/yag-config"
)

// Setup sets up the global logger based on the provided logging configuration.
func Setup(cfg Config) {
	out := os.Stderr
	var handler slog.Handler
	handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: cfg.Level})
	if cfg.Pretty {
		handler = tint.NewHandler(out, &tint.Options{Level: cfg.Level})
	}
	slog.SetDefault(slog.New(handler))
}

// Config represents the logging configuration.
type Config struct {
	Level  slog.Level
	Pretty bool
}

// Register registers the logging configuration flags with the provided parser.
func (cfg *Config) Register(y *yag.Parser) {
	y.Bool(&cfg.Pretty, "log_pretty", "Enable pretty colorful log output instead of default JSON")
	y.Value(logLevelFlag(&cfg.Level), "log_level", "Logging level. One of: DEBUG, INFO, WARN, ERROR")
}

// LogValue implements the [slog.LogValuer] interface.
func (cfg Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("level", cfg.Level.String()),
		slog.Bool("pretty", cfg.Pretty),
	)
}

func logLevelFlag(dest *slog.Level) *logLevelValue {
	return &logLevelValue{level: dest}
}

type logLevelValue struct {
	level *slog.Level
}

func (f *logLevelValue) String() string {
	return f.level.String()
}

func (f *logLevelValue) Set(value string) error {
	err := f.level.UnmarshalText([]byte(strings.ToUpper(value)))
	if err != nil {
		return err
	}
	return nil
}
