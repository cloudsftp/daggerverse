package main

import (
	"context"
	"fmt"

	"dagger/rust/internal/dagger"
)

// Lint the rust code
func (m *Rust) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	pkg string,
) error {
	lintCommand := []string{"cargo", "clippy"}

	if len(pkg) > 0 {
		lintCommand = append(lintCommand, "-p", pkg)
	} else {
		lintCommand = append(lintCommand, "--workspace")
	}

	lintCommand = append(lintCommand, "--", "-D", "warnings")

	if _, err := m.Builder(source).
		WithExec(lintCommand).
		Sync(ctx); err != nil {
		return fmt.Errorf("failing lints: %w", err)
	}

	return nil
}
