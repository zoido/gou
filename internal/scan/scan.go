package scan

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"slices"
	"sync"

	"mvdan.cc/xurls/v2"

	"github.com/zoido/gou/internal/model"
)

var urlRe = xurls.Relaxed()

// FindURLs scans trough [r] and returns found URLs as a slice of [model.Finding].
func FindURLs(ctx context.Context, r io.Reader) ([]model.Finding, error) {
	lines := make(chan []byte)
	errs := make(chan error)

	var wg sync.WaitGroup

	wg.Go(func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			select {
			case lines <- scanner.Bytes():
			case <-ctx.Done():
				return
			}
		}
		select {
		case <-ctx.Done():
		case errs <- scanner.Err():
		}
	})

	findings := make([]model.Finding, 0, 10) // Preallocate some.
	defer wg.Wait()

	for {
		select {
		case line := <-lines:
			matches := urlRe.FindAll(line, -1)
			for _, m := range matches {
				f, err := model.Builder{
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
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errs:
			if err != nil {
				return nil, err
			}
			return slices.Clip(findings), nil
		}
	}
}
