package db

import (
	"context"
	"log"
	// "time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "10.10.0.1:6379", // Change this to your Redis server address
		DB:   0,                // Default DB number
		// DialTimeout:  10 * time.Second,
		// ReadTimeout:  10 * time.Second,
		// MinIdleConns: 10,
		PoolSize: 1000,
	})
}

type RedisDB struct {
	rc *redis.Client
}

func (rdb *RedisDB) Getrc() *redis.Client {
	return rdb.rc
}

// KeyValueStoreDB
func (rdb *RedisDB) Get(k string) (string, error) {
	result, err := rdb.rc.Get(context.Background(), k).Result()
	if err == redis.Nil {
		return "0", nil
	} else if err != nil {
		log.Printf("Error fetching %s from Redis: %v", k, err)
		return "", err
	}
	return result, nil
}
func (rdb *RedisDB) Set(k string, v string) error {
	err := rdb.rc.Set(context.Background(), k, v, 0).Err()
	if err != nil {
		log.Printf("Error updating %s in Redis: %v", k, err)
		return err
	}
	return nil
}

func ConsRedisDB() *RedisDB {
	rc := InitRedis()
	return &RedisDB{rc: rc}
}
