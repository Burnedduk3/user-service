package ports

import (
	"context"
	"database/sql"
	"time"
)

// Database transaction interface
type Transactor interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

// Cache interface
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	HSet(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HExists(ctx context.Context, key, field string) (bool, error)
}

// Message queue interfaces
type MessagePublisher interface {
	Publish(ctx context.Context, exchange, routingKey string, message interface{}) error
}

type MessageConsumer interface {
	Consume(ctx context.Context, queue string, handler MessageHandler) error
}

type MessageHandler func(ctx context.Context, message []byte) error

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}
