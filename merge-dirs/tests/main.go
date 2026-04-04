package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type MergeDirsTests struct{}

type MergeDirectoriesTestCase struct {
	name     string
	dirs     []*dagger.Directory
	strategy dagger.MergeDirsMergeConflictStrategy
	expected *dagger.Directory
}

type MergeDirectoriesErrorTestCase struct {
	name     string
	dirs     []*dagger.Directory
	strategy dagger.MergeDirsMergeConflictStrategy
}

func (m *MergeDirsTests) All(ctx context.Context) error {
	if err := m.runSuccessTests(ctx); err != nil {
		return fmt.Errorf("success test cases failed: %w", err)
	}

	if err := m.runErrorTests(ctx); err != nil {
		return fmt.Errorf("error test cases failed: %w", err)
	}

	return nil
}

func (m *MergeDirsTests) runSuccessTests(ctx context.Context) error {
	tests := []MergeDirectoriesTestCase{
		// Conflict strategy ERROR
		{
			name:     "disjunct files at root",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a", "b": "b"}),
				buildDirectory(map[string]any{"c": "c", "d": "d"}),
			},
			expected: buildDirectory(map[string]any{"a": "a", "b": "b", "c": "c", "d": "d"}),
		},
		{
			name:     "disjunct files in subdir",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"sub": map[string]any{"a": "a"}}),
				buildDirectory(map[string]any{"sub": map[string]any{"b": "b"}}),
			},
			expected: buildDirectory(map[string]any{"sub": map[string]any{"a": "a", "b": "b"}}),
		},
		{
			name:     "disjunct subdirs",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"d1": map[string]any{"f1": "f1"}}),
				buildDirectory(map[string]any{"d2": map[string]any{"f2": "f2"}}),
			},
			expected: buildDirectory(map[string]any{
				"d1": map[string]any{"f1": "f1"},
				"d2": map[string]any{"f2": "f2"},
			}),
		},
		{
			name:     "nested disjunct",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{
							"c": map[string]any{"x": "x"},
						},
					},
				}),
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"y": "y"},
					},
				}),
			},
			expected: buildDirectory(map[string]any{
				"a": map[string]any{"b": map[string]any{
					"c": map[string]any{"x": "x"},
					"y": "y",
				}},
			}),
		},
		{
			name:     "deep nested",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"l1": map[string]any{
						"l2": map[string]any{
							"l3": map[string]any{"f1": "f1"},
						},
					},
				}),
				buildDirectory(map[string]any{
					"l1": map[string]any{
						"l2": map[string]any{
							"l3":  map[string]any{"f2": "f2"},
							"sib": "sib",
						},
					},
				}),
			},
			expected: buildDirectory(map[string]any{
				"l1": map[string]any{
					"l2": map[string]any{
						"l3":  map[string]any{"f1": "f1", "f2": "f2"},
						"sib": "sib",
					},
				},
			}),
		},
		{
			name:     "three dirs",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a"}),
				buildDirectory(map[string]any{"b": "b"}),
				buildDirectory(map[string]any{"c": "c"}),
			},
			expected: buildDirectory(map[string]any{"a": "a", "b": "b", "c": "c"}),
		},
		{
			name:     "four dirs shared subdir",
			strategy: dagger.MergeDirsMergeConflictStrategyErrorOnConflict,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"shared": map[string]any{"a": "a"},
				}),
				buildDirectory(map[string]any{
					"shared": map[string]any{"b": "b"},
				}),
				buildDirectory(map[string]any{
					"shared": map[string]any{"c": "c"},
				}),
				buildDirectory(map[string]any{
					"shared": map[string]any{"d": "d"},
					"other":  map[string]any{"x": "x"},
				}),
			},
			expected: buildDirectory(map[string]any{
				"shared": map[string]any{"a": "a", "b": "b", "c": "c", "d": "d"},
				"other":  map[string]any{"x": "x"},
			}),
		},

		// Conflict strategy KEEP_LEFT
		{
			name:     "keep left: root conflict",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "L"}),
				buildDirectory(map[string]any{"f": "R"}),
			},
			expected: buildDirectory(map[string]any{"f": "L"}),
		},
		{
			name:     "keep left: nested conflict",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"f": "L"},
					},
				}),
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"f": "R"},
					},
				}),
			},
			expected: buildDirectory(map[string]any{
				"a": map[string]any{
					"b": map[string]any{"f": "L"},
				},
			}),
		},
		{
			name:     "keep left: file over dir",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"x": "file"}),
				buildDirectory(map[string]any{"x": map[string]any{"inner": "i"}}),
			},
			expected: buildDirectory(map[string]any{"x": "file"}),
		},
		{
			name:     "keep left: three-way",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "1"}),
				buildDirectory(map[string]any{"f": "2"}),
				buildDirectory(map[string]any{"f": "3"}),
			},
			expected: buildDirectory(map[string]any{"f": "1"}),
		},

		// Conflict strategy KEEP_RIGHT
		{
			name:     "keep right: root conflict",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "L"}),
				buildDirectory(map[string]any{"f": "R"}),
			},
			expected: buildDirectory(map[string]any{"f": "R"}),
		},
		{
			name:     "keep right: nested conflict",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"f": "L"},
					},
				}),
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"f": "R"},
					},
				}),
			},
			expected: buildDirectory(map[string]any{
				"a": map[string]any{
					"b": map[string]any{"f": "R"},
				},
			}),
		},
		{
			name:     "keep right: dir replaces file",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"x": "file"}),
				buildDirectory(map[string]any{"x": map[string]any{"inner": "i"}}),
			},
			expected: buildDirectory(map[string]any{"x": map[string]any{"inner": "i"}}),
		},
		{
			name:     "keep right: three-way",
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "1"}),
				buildDirectory(map[string]any{"f": "2"}),
				buildDirectory(map[string]any{"f": "3"}),
			},
			expected: buildDirectory(map[string]any{"f": "3"}),
		},
	}

	for _, test := range tests {
		err := test.run(ctx)
		if err != nil {
			return fmt.Errorf("test case '%s' failed: %w", test.name, err)
		}
	}

	return nil
}

func (m *MergeDirsTests) runErrorTests(ctx context.Context) error {
	tests := []MergeDirectoriesErrorTestCase{
		{
			name: "file conflict at root",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "L"}),
				buildDirectory(map[string]any{"f": "R"}),
			},
		},
		{
			name: "file conflict in subdir",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"sub": map[string]any{"f": "L"}}),
				buildDirectory(map[string]any{"sub": map[string]any{"f": "R"}}),
			},
		},
		{
			name: "file conflict nested deep",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{
							"c": map[string]any{"f": "L"},
						},
					},
				}),
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{
							"c": map[string]any{"f": "R"},
						},
					},
				}),
			},
		},
		{
			name: "type mismatch: dir then file",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"x": map[string]any{"i": "i"}}),
				buildDirectory(map[string]any{"x": "file"}),
			},
		},
		{
			name: "type mismatch: file then dir",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"x": "file"}),
				buildDirectory(map[string]any{"x": map[string]any{"i": "i"}}),
			},
		},
		{
			name: "type mismatch nested",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"x": "file"},
					},
				}),
				buildDirectory(map[string]any{
					"a": map[string]any{
						"b": map[string]any{"x": map[string]any{"i": "i"}},
					},
				}),
			},
		},
		{
			name: "conflict in three-way",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"f": "1"}),
				buildDirectory(map[string]any{"other": "x"}),
				buildDirectory(map[string]any{"f": "2"}),
			},
		},
		{
			name: "conflict early in four-way",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"x": "a"}),
				buildDirectory(map[string]any{"y": "b"}),
				buildDirectory(map[string]any{"z": "c"}),
				buildDirectory(map[string]any{"x": "d"}),
			},
		},
	}

	for _, test := range tests {
		err := test.run(ctx)
		if err != nil {
			return fmt.Errorf("test case '%s' failed: %w", test.name, err)
		}
	}

	return nil
}

func (t *MergeDirectoriesTestCase) run(ctx context.Context) error {
	merged := dag.MergeDirs().Merge(
		t.dirs,
		dagger.MergeDirsMergeOpts{
			Strategy: t.strategy,
		},
	)

	return assertDirectory(ctx, t.expected, merged)
}

func (t *MergeDirectoriesErrorTestCase) run(ctx context.Context) error {
	merged := dag.MergeDirs().Merge(
		t.dirs,
		dagger.MergeDirsMergeOpts{
			Strategy: t.strategy,
		},
	)

	_, err := merged.Entries(ctx, dagger.DirectoryEntriesOpts{})
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}

	return nil
}
