package cache

import (
	"time"
	"gopkg.in/redis.v5"
	"strconv"
)

type CRedis struct {
	client *redis.Client
}

func newCacheRedis(host, database string) (*CRedis, error) {
	db, err := strconv.Atoi(database)
	if err != nil {
		return nil, err
	}
	c := new(CRedis)
	c.client = redis.NewClient(&redis.Options{
		Addr:     host,
		DB:       db,
	})

	return c, nil
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