package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func New(url string) (*Redis, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis url: %w", err)
	}

	var rd Redis

	rd.Client = redis.NewClient(opt)

	return &rd, nil
}

func (r *Redis) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
