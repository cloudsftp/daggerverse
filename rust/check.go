package main

import (
	"context"
	"fmt"

	"dagger/rust/internal/dagger"
)

// Check the rust code
func (m *Rust) Check(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	name string,
) error {
	checkCommand := []string{"cargo", "check"}

	if len(name) > 0 {
		checkCommand = append(checkCommand, "-p", name)
	}

	if _, err := m.builder(source).
		WithExec(checkCommand).
		Sync(ctx); err != nil {
		return fmt.Errorf("unsuccessful check of source code: %w", err)
	}

	return nil
}
