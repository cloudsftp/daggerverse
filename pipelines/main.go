package main

import (
	"context"
	"fmt"
)

type CloudsDaggerModules struct{}

// Run Tests
func (m *CloudsDaggerModules) Test(ctx context.Context) error {
	var err error

	err = dag.MergeDirsTests().All(ctx)
	if err != nil {
		return fmt.Errorf("merge directories tests failed: %w", err)
	}

	return nil
}
