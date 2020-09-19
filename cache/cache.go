package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cache := &Cache{
		client: client,
	}

	return cache
}

func (c *Cache) Get(key string) ([]byte, error) {
	return c.client.Get(key).Bytes()
}

func (c *Cache) Set(key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshalling the data to JSON ", err)
	}
	return c.client.Set(key, jsonData, expiration).Err()
}
