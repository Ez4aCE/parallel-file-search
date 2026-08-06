package search

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func createBenchmarkFiles(b *testing.B, fileCount int) []string {
	b.Helper()
	dir := b.TempDir()
	paths := make([]string, 0, fileCount)
	for i := 0; i < fileCount; i++ {
		path := filepath.Join(dir, fmt.Sprintf("log_%d.log", i))
		content := strings.Repeat(
			"INFO startup\nERROR timeout\nINFO request\n",
			1000,
		)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
		paths = append(paths, path)

	}
	return paths
}

func BenchmarkFilesSearchSequential(b *testing.B) {
	paths := createBenchmarkFiles(b, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FilesSearchSequential(paths, "ERROR")
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkFilesSearchConcurrent(b *testing.B) {
	paths := createBenchmarkFiles(b, b.N)
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := FilesSearchConcurrent(ctx, paths, "ERROR", runtime.NumCPU())
		if err != nil {
			b.Fatal(err)
		}
	}
}
