package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/Ez4aCE/parallel-file-search/internal/search"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt)
	defer stop()
	term := flag.String("term", "", "target term")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	flag.Parse()
	args := flag.Args()

	start := time.Now()

	result, err := search.FilesSearchConcurrent(ctx, args, *term, *workers)
	duration := time.Since(start)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	printResult(result)
	fmt.Println("Completed in ", duration.Milliseconds(), "ms")

}
func printResult(result search.Result) {
	totalFiles := len(result.Matches) + len(result.Errors)
	fmt.Printf("Searched %d\n", totalFiles)
	for path, matches := range result.Matches {
		fmt.Println(path)

		if len(matches) == 0 {
			fmt.Println("  No matches")
		} else {
			for _, line := range matches {
				fmt.Println(" ", line)
			}
		}

		fmt.Println()
	}
	if len(result.Errors) > 0 {
		fmt.Println("Errors")
		fmt.Println("------")

		for path, err := range result.Errors {
			fmt.Printf("%s: %v\n", path, err)
		}
	}
}
