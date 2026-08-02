package exec

import (
	"context"

	"github.com/zoido/gou/internal/exec/app"
)

// Exec is the entry point for the gou command.
func Exec(ctx context.Context, args []string) int {
	err := app.Run(ctx, args)
	if err != nil {
		return 1
	}
	return 0
}
