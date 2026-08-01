package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Ez4aCE/parallel-file-search/internal/search"
)

func main() {
	term := flag.String("term", "", "target term")
	flag.Parse()
	args := flag.Args()

	start := time.Now()
	result, err := search.FilesSearchConcurrent(args, *term)
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
	for k, v := range result.Matches {
		fmt.Println(k, v)
	}
	for k, v := range result.Errors {
		fmt.Println(k, v)
	}
}
