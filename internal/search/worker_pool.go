package search

import (
	"context"
	"sync"
)

func worker(
	ctx context.Context,
	jobs <-chan string,
	results chan<- FileResult,
	wg *sync.WaitGroup,
	term string,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-jobs:
			if !ok {
				return
			}
			matches, err := SingleFileSearch(path, term)
			select {
			case <-ctx.Done():
				return

			case results <- FileResult{
				Path:    path,
				Matches: matches,
				Err:     err,
			}:
			}

		}
	}
}
