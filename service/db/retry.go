package db

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

func RetryOperation(operation func() error, maxRetries int) error {
	var err error
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = operation()

		if err == nil {
			return nil
		}

		// Check if error is retryable
		if isRetryableError(err) {
			log.Printf("Database operation failed, retrying in %v (attempt %d/%d): %v",
				backoff, attempt+1, maxRetries, err)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
			continue
		}

		// Non-retryable error
		return err
	}

	return err
}

// isRetryableError determines if a database error can be retried
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if err == sql.ErrConnDone {
		return true
	}

	if pqErr, ok := err.(interface{ sqlState() string }); ok {
		code := pqErr.sqlState()
		return strings.HasPrefix(code, "08") ||
			strings.HasPrefix(code, "57") ||
			strings.HasPrefix(code, "53")
	}
	errStr := err.Error()
	retryableErrors := []string{
		"connection reset by peer",
		"broken pipe",
		"connection refused",
		"no connection to the server",
		"unexpected EOF",
		"connection timed out",
	}

	for _, msg := range retryableErrors {
		if strings.Contains(errStr, msg) {
			return true
		}
	}

	return false
}
