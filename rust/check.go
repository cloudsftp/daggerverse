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
	pkg string,
) error {
	checkCommand := []string{"cargo", "check"}

	if len(pkg) > 0 {
		checkCommand = append(checkCommand, "-p", pkg)
	} else {
		checkCommand = append(checkCommand, "--workspace")
	}

	if _, err := m.builder(source).
		WithExec(checkCommand).
		Sync(ctx); err != nil {
		return fmt.Errorf("unsuccessful check of source code: %w", err)
	}

	return nil
}
