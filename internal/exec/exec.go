package exec

import (
	"context"
)

// Exec is the entry point for the gou command.
func Exec(ctx context.Context, args []string) int {
	_ = ctx
	_ = args
	return 0
}
