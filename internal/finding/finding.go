package model

import (
	"errors"
	"fmt"
	"strings"
)

// Finding represents one item. An URl found in the text and optional context.
type Finding struct {
	valid bool

	url string
}

func (f Finding) URL() string {
	f.ensureValid()
	return f.url
}

func (f Finding) ensureValid() {
	if !f.valid {
		panic("finding was not initialized properly")
	}
}

// Builder is a helper struct for initialising [Finding].
type Builder struct {
	// Prevents comparability and use of unkeyed literals for the builder.
	// Borrowed from the gRPC opaque API.
	_ [0]func()

	URL string
}

func (b Builder) Build() (Finding, error) {
	var errs []error

	trimmedURL := strings.TrimSpace(b.URL)
	if trimmedURL == "" {
		errs = append(errs, fmt.Errorf("url cannot be empty"))
	}

	if err := errors.Join(errs...); err != nil {
		return Finding{}, err
	}

	return Finding{
		url: trimmedURL,

		valid: true,
	}, nil
}
