package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	client *redis.Client
	prefix string
	group  singleflight.Group
}

func NewCache(client *redis.Client, prefix string) *Cache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

func (c *Cache) key(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	if c.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	data, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(key), data, expiration).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Del(ctx, c.key(key)).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) bool {
	if c.client == nil {
		return false
	}
	n, err := c.client.Exists(ctx, c.key(key)).Result()
	return err == nil && n > 0
}

func (c *Cache) GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	if c.client != nil {
		err := c.Get(ctx, key, dest)
		if err == nil {
			return nil
		}
	}

	result, err, _ := c.group.Do(key, func() (interface{}, error) {
		data, err := fn()
		if err != nil {
			return nil, err
		}

		if c.client != nil {
			if setErr := c.Set(ctx, key, data, expiration); setErr != nil {
				return nil, setErr
			}
		}

		return data, nil
	})

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}
