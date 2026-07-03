// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package parproc

import (
	"iter"
	"runtime"
	"sync"
)

type Result[P, R any] struct {
	Payload P
	Outcome R
	Error   error
}

func ParProc[P, R any](
	proc func(P) (R, error),
	payloads []P,
	nworkers int,
) iter.Seq[Result[P, R]] {
	return func(yield func(Result[P, R]) bool) {
		out := make(chan Result[P, R], len(payloads))
		done := make(chan struct{})

		// If the iterator exits for any reason (normal finish or early break),
		// closing 'done' broadcasts a signal to stop processing immediately.
		defer close(done)

		go func() {
			defer close(out)
			if nworkers <= 0 {
				nworkers = runtime.GOMAXPROCS(0)
			}
			var wg sync.WaitGroup
			sem := make(chan int, nworkers)

			for i, path := range payloads {
				// Check cancellation before pulling a new worker slot
				select {
				case <-done:
					return
				default:
				}

				select {
				case <-done:
					return
				case sem <- i:
				}

				wg.Add(1)
				go func(p P) {
					defer wg.Done()
					defer func() { <-sem }()

					// Check cancellation right before doing heavy lifting
					select {
					case <-done:
						return
					default:
					}

					result, err := proc(p)

					// Safely deliver the result
					out <- Result[P, R]{Payload: p, Outcome: result, Error: err}
				}(path)
			}
			wg.Wait()
		}()

		for r := range out {
			if !yield(r) {
				break // Triggers 'defer close(done)', tearing down the pool safely
			}
		}
	}
}
