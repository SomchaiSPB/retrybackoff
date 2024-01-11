package retrybackoff

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func mockOperation(retryUntilSuccess int, timestamps *[]time.Time) func() error {
	counter := 0
	return func() error {
		*timestamps = append(*timestamps, time.Now())
		if counter < retryUntilSuccess {
			counter++
			return fmt.Errorf("error")
		}
		return nil
	}
}

func TestRetrier(t *testing.T) {
	type testCase struct {
		name            string
		operation       func() error
		maxRetries      int
		maxBackoff      time.Duration
		wantErr         bool
		expectedRetries int
		timestamps      []time.Time
	}

	tests := []testCase{
		{
			name:            "Success on first try",
			operation:       nil,
			maxRetries:      3,
			maxBackoff:      2 * time.Second,
			wantErr:         false,
			expectedRetries: 0,
		},
		{
			name:            "Success after retries",
			operation:       nil,
			maxRetries:      3,
			maxBackoff:      2 * time.Second,
			wantErr:         false,
			expectedRetries: 2,
		},
		{
			name:            "Error after retries",
			operation:       nil,
			maxRetries:      2,
			maxBackoff:      2 * time.Second,
			wantErr:         true,
			expectedRetries: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.operation = mockOperation(tc.expectedRetries, &tc.timestamps)

			err := RetryWithBackoff(tc.operation, tc.maxRetries, tc.maxBackoff)

			if (tc.wantErr && err == nil) || (!tc.wantErr && err != nil) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}

			for i := 1; i < len(tc.timestamps); i++ {
				actualBackoff := tc.timestamps[i].Sub(tc.timestamps[i-1])
				expectedBackoff := time.Duration(math.Pow(2, float64(i-1))) * time.Second
				if expectedBackoff > tc.maxBackoff {
					expectedBackoff = tc.maxBackoff
				}

				if actualBackoff < expectedBackoff {
					t.Errorf("Backoff time too short. Retry %d: expected at least %v, got %v", i, expectedBackoff, actualBackoff)
				}
			}
		})
	}
}
