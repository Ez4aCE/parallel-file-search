package search

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

type searchFilesTestCase struct {
	name            string
	files           map[string]string
	missingFiles    []string
	searchTerm      string
	expectedMatches map[string][]string
	expectedErrors  []string
	expectErr       error
}

func TestSearchFiles(t *testing.T) {
	tests := []searchFilesTestCase{
		{
			name: "multiple_files_success",
			files: map[string]string{
				"app.log": `INFO startup
ERROR timeout`,
				"api.log":    `ERROR invalid token`,
				"worker.log": `INFO started`,
			},
			searchTerm: "ERROR",
			expectedMatches: map[string][]string{
				"app.log": {
					"ERROR timeout",
				},
				"api.log": {
					"ERROR invalid token",
				},
				"worker.log": {},
			},
		},
		{
			name: "partial_failure",
			files: map[string]string{
				"app.log": `ERROR timeout`,
			},
			missingFiles: []string{
				"missing.log",
			},
			searchTerm: "ERROR",
			expectedMatches: map[string][]string{
				"app.log": {
					"ERROR timeout",
				},
			},
			expectedErrors: []string{
				"missing.log",
			},
		},
		{
			name:       "no_files_provided",
			searchTerm: "ERROR",
			expectErr:  ErrNoFilesProvided,
		},
		{
			name: "empty_search_term",
			files: map[string]string{
				"app.log": `ERROR timeout`,
			},
			expectErr: ErrEmptySearchTerm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			var paths []string

			for name, content := range tt.files {
				path := filepath.Join(dir, name)

				err := os.WriteFile(
					path,
					[]byte(content),
					0644,
				)
				if err != nil {
					t.Fatal(err)
				}

				paths = append(paths, path)
			}

			for _, name := range tt.missingFiles {
				path := filepath.Join(dir, name)
				paths = append(paths, path)
			}

			result, err := FilesSearch(paths, tt.searchTerm)

			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("got error %v, want %v", err, tt.expectErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for path, expected := range tt.expectedMatches {
				fullPath := filepath.Join(dir, path)

				got, ok := result.Matches[fullPath]
				if !ok {
					t.Fatalf("missing result for %s", path)
				}

				if !slices.Equal(got, expected) {
					t.Errorf(
						"path=%s got=%v want=%v",
						path,
						got,
						expected,
					)
				}
			}

			for _, file := range tt.expectedErrors {
				fullPath := filepath.Join(dir, file)

				if _, ok := result.Errors[fullPath]; !ok {
					t.Errorf(
						"expected error for file %s",
						file,
					)
				}
			}
		})
	}
}
