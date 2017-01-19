package cache

import "time"

type CStorage interface {
	Set(key string, data string, expr time.Duration) error
	Get(key string) string
}
