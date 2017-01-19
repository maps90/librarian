package cache

import (
	"time"
	"gopkg.in/redis.v5"
)

type CRedis struct {
	client *redis.Client
}

func NewCRedis(client *redis.Client) *CRedis {
	c := new(CRedis)
	c.client = client

	return c
}

func (c *CRedis) Set(key string, data string, expr time.Duration) error {
	if err := c.client.Set(key, data, expr).Err(); err != nil {
		return err
	}

	return nil
}

func (c *CRedis) Get(key string) string {
	b, err := c.client.Get(key).Result()
	if err != nil {
		return ""
	}

	return b
}