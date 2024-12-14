package redis

import (
	"context"
	"time"
)

//go:generate mockgen -source=ports.go -destination=service_redis_mock_test.go -package=service
type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
