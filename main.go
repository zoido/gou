package main

import (
	"context"
	"os"

	"github.com/zoido/gou/internal/exec"
)

func main() {
	os.Exit(exec.Exec(context.Background(), os.Args[1:]))
}
