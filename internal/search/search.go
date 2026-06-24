package search

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Result struct {
	Matches map[string][]string
	Errors  map[string]error
}

var ErrNoFilesProvided = errors.New("no files provided")
var ErrEmptySearchTerm = errors.New("search term is empty")

func FilesSearch(paths []string, term string) (Result, error) {

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
			continue
		}
		result.Matches[path] = matches
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
