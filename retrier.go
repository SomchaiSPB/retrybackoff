package retrybackoff

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"time"
)

var ErrMaxRetriesExceeded = errors.New("reached max retries without success")

func RetryWithBackoff(operation func() error, maxRetries int, maxBackoff time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err != nil {
			if i == maxRetries-1 {
				log.Printf("max connection retries exceeded: %v", err)
				return err
			}

			backoffTime := time.Duration(math.Pow(2, float64(i))) * time.Second

			if backoffTime > maxBackoff {
				backoffTime = maxBackoff
			}

			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			backoffTimeWithJitter := backoffTime + jitter

			log.Printf("retry %d: reconnecting in %v...", i+1, backoffTimeWithJitter)
			time.Sleep(backoffTimeWithJitter)

			continue
		}

		return nil
	}

	return ErrMaxRetriesExceeded
}
