package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/config"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Client *redis.Client

// InitRedis connects to Redis and sets up the global client
func InitRedis(ctx context.Context) error {
	if config.AppConfig == nil || config.AppConfig.Redis.URL == "" {
		return fmt.Errorf("redis URL not configured")
	}

	opts, err := redis.ParseURL(config.AppConfig.Redis.URL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Optimize pool settings
	opts.PoolSize = 100
	opts.MinIdleConns = 10
	opts.ConnMaxLifetime = 5 * time.Minute

	Client = redis.NewClient(opts)

	// Ping to verify connection
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Log.Info("Connected to Redis successfully", zap.String("url", opts.Addr))
	return nil
}

// Close gracefully closes the Redis connection pool
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
