package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxRetries = 5
	defaultRetryDelay = 2 * time.Second
)

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	return ConnectWithRetry(ctx, dbURL, defaultMaxRetries, defaultRetryDelay)
}

func ConnectWithRetry(ctx context.Context, dbURL string, maxRetries int, retryDelay time.Duration) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pool, lastErr = pgxpool.New(ctx, dbURL)
		if lastErr != nil {
			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}
			continue
		}

		if pingErr := pool.Ping(ctx); pingErr != nil {
			pool.Close()
			lastErr = pingErr
			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}
			continue
		}

		return pool, nil
	}

	return nil, fmt.Errorf("connect to database after %d attempts: %w", maxRetries, lastErr)
}

func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("database pool is nil")
	}
	return pool.Ping(ctx)
}
