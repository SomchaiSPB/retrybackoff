# GoRetryBackoff

GoRetryBackoff is a Go package providing a robust retry mechanism with exponential backoff and jitter. It's designed to handle transient failures in distributed systems, enhancing stability and performance in Go applications.

## Features

- **Exponential Backoff:** Increases the delay between retries exponentially.
- **Jitter:** Adds randomness to the backoff intervals to prevent thundering herd problem.
- **Customizable:** Easily configure max retries and max backoff duration.

## Installation

```bash
go get github.com/SomchaiSPB/retrybackoff
```

## Usage

```go
package main

import (
    retrier "github.com/SomchaiSPB/retrybackoff"
    "log"
    "net/http"
    "time"
)

func main() {
    operation := func() error {
        // Replace with your operation
        _, err := http.Get("http://example.com")
        return err
    }

    err := retrier.RetryWithBackoff(operation, 5, 2*time.Minute)
    if err != nil {
        log.Fatalf("Operation failed: %v", err)
    }
}
```