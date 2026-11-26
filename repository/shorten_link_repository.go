package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// const expirationTime = 1 * time.Hour // 1 hour

type URLStorage interface {
	SaveURL(ctx context.Context, key, url string, expirationTime time.Duration) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

func (r *urlStorage) SaveURL(ctx context.Context, key, url string, expirationTime time.Duration) error {
	return r.c.Set(ctx, key, url, expirationTime).Err()
}

func (r *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	url, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}
