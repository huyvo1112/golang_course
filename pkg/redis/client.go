package redis

import goredis "github.com/redis/go-redis/v9"

// Create new redis client
func NewClient(envPrefix string) (*goredis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb, nil
}
