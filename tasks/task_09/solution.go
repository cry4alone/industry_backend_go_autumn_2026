package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return []Result[R]{}, nil
	}
	if workers > len(in) {
		workers = len(in)
	}

	results := make([]Result[R], len(in))
	var next int
	var nextMu sync.Mutex
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for {
				if ctx.Err() != nil {
					return
				}

				nextMu.Lock()
				if ctx.Err() != nil || next == len(in) {
					nextMu.Unlock()
					return
				}
				i := next
				next++
				nextMu.Unlock()

				if ctx.Err() != nil {
					return
				}
				value, err := fn(ctx, in[i])
				if err != nil {
					results[i].Err = err
					continue
				}
				results[i].Value = value
			}
		})
	}

	wg.Wait()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
