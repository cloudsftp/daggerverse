package main

import (
	"context"
	"fmt"

	"dagger/rust/internal/dagger"
)

// Test the rust code
func (m *Rust) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	pkg string,
) error {
	testCommand := []string{"cargo", "test"}

	if len(pkg) > 0 {
		testCommand = append(testCommand, "-p", pkg)
	}

	if _, err := m.builder(source).
		WithExec(testCommand).
		Sync(ctx); err != nil {
		return fmt.Errorf("failing tests: %w", err)
	}

	return nil
}
