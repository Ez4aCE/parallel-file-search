package search

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"
	"sync"
)

type Result struct {
	Matches map[string][]string
	Errors  map[string]error
}

type FileResult struct {
	Path    string
	Matches []string
	Err     error
}

var ErrNoFilesProvided = errors.New("no files provided")
var ErrEmptySearchTerm = errors.New("search term is empty")

func FilesSearchConcurrent(ctx context.Context, paths []string, term string, workers int) (Result, error) {

	if len(paths) == 0 {
		return Result{}, ErrNoFilesProvided
	}
	if len(term) == 0 {
		return Result{}, ErrEmptySearchTerm
	}

	result := Result{
		Matches: make(map[string][]string),
		Errors:  make(map[string]error),
	}
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan string)
	results := make(chan FileResult)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, results, &wg, term)
	}
	go func() {
		defer close(jobs)
		for _, path := range paths {
			jobs <- path
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Err != nil {
			result.Errors[res.Path] = res.Err
			continue
		}
		result.Matches[res.Path] = res.Matches
	}

	return result, nil
}
func FilesSearchSequential(paths []string, term string) (Result, error) {
	if len(paths) == 0 {
		return Result{}, ErrNoFilesProvided
	}
	if len(term) == 0 {
		return Result{}, ErrEmptySearchTerm
	}
	result := Result{
		Matches: make(map[string][]string),
		Errors:  make(map[string]error),
	}
	for _, path := range paths {
		matches, err := SingleFileSearch(path, term)
		if err != nil {
			result.Errors[path] = err
		} else {
			result.Matches[path] = matches
		}
	}
	return result, nil
}
func SingleFileSearch(path string, term string) ([]string, error) {
	var matches []string
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, term) {
			matches = append(matches, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}
