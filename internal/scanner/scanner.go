package scanner

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/zoido/gou/internal/domain/finding"
	"mvdan.cc/xurls/v2"
)

var urlRe = xurls.Relaxed()

func FindURLs(ctx context.Context, r io.Reader) ([]finding.Finding, error) {
	lines := make(chan []byte)
	defer close(lines)
	errs := make(chan error)
	defer close(errs)

	var wg sync.WaitGroup

	wg.Go(func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines <- scanner.Bytes()
		}
		errs <- scanner.Err()
	})

	findings := make([]finding.Finding, 0, 10) // Preallocate some.

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case err := <-errs:
			if err != nil {
				return nil, err
			}
			wg.Wait()
			return slices.Clip(findings), nil
		case line := <-lines:
			matches := urlRe.FindAll(line, -1)
			for _, m := range matches {
				f, err := finding.Builder{
					URL: string(m),
				}.Build()
				if err != nil {
					// This should not happen. But if it does we lean towards returning all that do
					// not fail.
					slog.Debug("Failed to initialize finding", "error", err)
					continue
				}
				findings = append(findings, f)
			}
		}
	}
}
