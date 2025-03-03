package utils

import (
	"context"
	"time"
)

func NewTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

const DefaultQueryTimeout = 30 * time.Second
