package db

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type BatchedRedisDB struct {
	rc *redis.Client
}

// todo: use generic operations
type Op interface{}
type GetOp struct {
	K string
}
type SetOp struct {
	K string
	V string
}

var batch []Op
var batchSize = 100

func AppendToBatch(op Op) {
	switch v := op.(type) {
	case GetOp:
		fmt.Println("Appending GetOp:", v.K)
		batch = append(batch, v) // Append the GetOp to the batch
	case SetOp:
		fmt.Println("Appending SetOp:", v.K, v.V)
		batch = append(batch, v) // Append the SetOp to the batch
	default:
		fmt.Println("Unknown operation type")
	}
}

func GetBatch() []Op {
	return batch
}

// KeyValueStoreDB
func (rdb *BatchedRedisDB) Get(k string) (string, error) {
	// Create an operation
	// Push operation on batch list processor list
	// Create a Task that waits on the result of operation
	// Wait an task and return the result
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

// This version uses batch for a single
func ConsBatchedRedisDB() *BatchedRedisDB {
	rc := initRedis()
	return &BatchedRedisDB{rc: rc}
}
