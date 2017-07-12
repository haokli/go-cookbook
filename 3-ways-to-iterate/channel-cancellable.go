package main

import (
	"context"
	"fmt"
	"log"
)

// To run:
// go run channel-cancellable.go

// IntWithError combines an integer value and an error
type IntWithError struct {
	Int int
	Err error
}

func generateEvenNumbers(ctx context.Context, max int) chan IntWithError {
	ch := make(chan IntWithError)
	go func() {
		defer close(ch)
		if max < 0 {
			ch <- IntWithError{
				Err: fmt.Errorf("'max' is %d and should be >= 0", max),
			}
			return
		}

		for i := 2; i <= max; i += 2 {
			if ctx != nil {
				// if context was cancelled, we stop early
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
			ch <- IntWithError{
				Int: i,
			}
		}
	}()
	return ch
}

func printEvenNumbersCancellable(max int, stopAt int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for val := range generateEvenNumbers(ctx, max) {
		if val.Err != nil {
			log.Fatalf("Error: %s\n", val.Err)
		}
		if val.Int > stopAt {
			// notice we keep going to drain the channel
			cancel()
		}
		fmt.Printf("%d\n", val.Int)
	}
}

func printEvenNumbersCancellable2(max int, stopAt int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	iter := generateEvenNumbers(ctx, max)
	for val := range iter {
		if val.Err != nil {
			log.Fatalf("Error: %s\n", val.Err)
		}
		if val.Int > stopAt {
			// notice we keep going to drain the channel
			cancel()
			break
		}
		fmt.Printf("%d\n", val.Int)
	}
	// drain the channel
	for _ = range iter {
	}
}

func main() {
	fmt.Printf("Even numbers up to 20, cancel at 8:\n")
	printEvenNumbersCancellable(20, 8)
	fmt.Printf("Even numbers up to 20, cancel at 8:\n")
	printEvenNumbersCancellable2(20, 8)
}
