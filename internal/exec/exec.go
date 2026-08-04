package exec

import (
	"context"
	"log/slog"

	"github.com/zoido/gou/internal/exec/app"
)

// Exec is the entry point for the gou command.
func Exec(ctx context.Context, args []string) int {
	err := app.Run(ctx, args)
	if err != nil {
		slog.Error("Application failed", "error", err)
		return 1
	}
	return 0
}
