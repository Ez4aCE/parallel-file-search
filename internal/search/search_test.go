package search

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

type testCase struct {
	name            string
	fileContent     string
	searchTerm      string
	expectedMatches []string
	expectErr       bool
}

func TestSearchFile(t *testing.T) {
	tests := []testCase{
		{
			name:        "multiple_matches",
			fileContent: "INFO startup\nERROR timeout\nINFO request\nERROR invalid token",
			searchTerm:  "ERROR",
			expectedMatches: []string{
				"ERROR timeout",
				"ERROR invalid token",
			},
			expectErr: false,
		},
		{
			name:            "no_matches",
			fileContent:     "INFO startup\nINFO request\nINFO done",
			searchTerm:      "ERROR",
			expectedMatches: []string{},
			expectErr:       false,
		},
		{
			name:            "empty_file",
			fileContent:     "",
			searchTerm:      "ERROR",
			expectedMatches: []string{},
			expectErr:       false,
		},
		{
			name:        "empty_search_term",
			fileContent: "INFO startup\nINFO request\nINFO done",
			searchTerm:  "",
			expectedMatches: []string{
				"INFO startup",
				"INFO request",
				"INFO done",
			},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "app.log")
			content := []byte(tt.fileContent)
			err := os.WriteFile(path, content, 0644)

			if err != nil {
				t.Fatal(err)
			}

			matches, err := SearchFile(path, tt.searchTerm)

			if !slices.Equal(matches, tt.expectedMatches) {
				t.Errorf("got %v, want %v", matches, tt.expectedMatches)
			}
			if tt.expectErr && err == nil {
				t.Errorf("got %v, want error", matches)
			}
		})
	}
}

func TestSearchFile_FileNotFound(t *testing.T) {

	matches, err := SearchFile("notHere.log", "ERROR")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if matches != nil {
		t.Errorf("got %v, want nil", matches)
	}
}
