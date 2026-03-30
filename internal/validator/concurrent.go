package validator

import (
	"runtime"
	"sort"
	"sync"

	"github.com/chuck/openspec-go/internal/model"
)

// ValidationFunc is a function that validates an item and returns issues.
type ValidationFunc func() (model.ValidationItem, error)

// ValidateConcurrent runs validation functions concurrently with bounded parallelism.
// Results are returned in deterministic order sorted by item ID.
func ValidateConcurrent(items []ValidationFunc, concurrency int) []model.ValidationItem {
	if len(items) == 0 {
		return nil
	}

	// Single item: run directly, no pool
	if len(items) == 1 {
		result, _ := items[0]()
		return []model.ValidationItem{result}
	}

	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}

	results := make([]model.ValidationItem, len(items))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for i, fn := range items {
		wg.Add(1)
		go func(idx int, f ValidationFunc) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			result, _ := f()
			results[idx] = result
		}(i, fn)
	}

	wg.Wait()

	// Sort by ID for deterministic output
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results
}
