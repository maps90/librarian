package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type CMemcached struct {
	client *memcache.Client
}

func NewCMemcached(client *memcache.Client) *CMemcached {
	c := new(CMemcached)
	c.client = client

	return c
}

func (c *CMemcached) Set(key string, data string, expr time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(data),
		Expiration: int32(time.Now().Add(expr).Unix()),
	}
	if err := c.client.Set(item); err != nil {
		return err
	}
	return nil
}

func (c *CMemcached) Get(key string) string {
	b, err := c.client.Get(key)
	if err != nil {
		return ""
	}
	return string(b.Value)
}
