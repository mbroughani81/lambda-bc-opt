package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type BatchedRedisDB struct {
	rc *redis.Client
}

// KeyValueStoreDB
func (rdb *BatchedRedisDB) Get(k string) (string, error) {
	result, err := rdb.rc.Get(context.Background(), k).Result()
	if err == redis.Nil {
		return "0", nil
	} else if err != nil {
		log.Printf("Error fetching %s from Redis: %v", k, err)
		return "", err
	}
	return result, nil
}
func (rdb *BatchedRedisDB) Set(k string, v string) error {
	err := rdb.rc.Set(context.Background(), k, v, 0).Err()
	if err != nil {
		log.Printf("Error updating %s in Redis: %v", k, err)
		return err
	}
	return nil
}

func ConsBatchedRedisDB() *RedisDB {
	rc := initRedis()
	return &RedisDB{rc: rc}
}
