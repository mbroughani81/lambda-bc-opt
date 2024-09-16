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

type BatchOp struct {
	op Op
	ch chan string
}

var batch []BatchOp
var batchSize = 3

func AppendToBatch(rdb *BatchedRedisDB, op Op, ch chan string) {
	go func() {
		switch v := op.(type) {
		case GetOp:
			fmt.Println("Appending GetOp:", v.K)
			batch = append(batch, BatchOp{op, ch}) // Append the GetOp to the batch
		case SetOp:
			fmt.Println("Appending SetOp:", v.K, v.V)
			batch = append(batch, BatchOp{op, ch}) // Append the SetOp to the batch
		default:
			log.Fatalln("Unknown operation type")
		}
		if len(batch) >= batchSize {
			ctx := context.Background()
			var redisResponses []redis.Cmder
			// do the operations in pipleline
			pipe := rdb.rc.Pipeline()
			for _, x := range batch {
				cur_op := x.op
				switch v := cur_op.(type) {
				case GetOp:
					result := pipe.Get(ctx, v.K)
					redisResponses = append(redisResponses, result)
				case SetOp:
					result := pipe.Set(ctx, v.K, v.V, 0)
					redisResponses = append(redisResponses, result)
				default:
					log.Fatalln("Unknown operation type")
				}
			}
			// executing the pipeline
			log.Printf("Executing the pipeline => %#v", batch)
			_, err := pipe.Exec(ctx)
			if err != nil && err != redis.Nil {
				log.Fatalf("Executing the pipeline failed! => %v", err)
			}
			if err == redis.Nil {
				log.Printf("REDIS IS NIL (some of the keys are not found)")
			}
			for index, resp := range redisResponses {
				log.Printf("redis-responses => %s\n", resp.String())
				op := batch[index].op
				ch := batch[index].ch
				switch op.(type) {
				case GetOp:
					if getResp, ok := resp.(*redis.StringCmd); ok {
						ch <- getResp.Val()
					} else {
						ch <- "error"
					}
				case SetOp:
					if setResp, ok := resp.(*redis.StatusCmd); ok {
						ch <- setResp.Val()
					} else {
						ch <- "error"
					}
				default:
					log.Fatalln("Unknown operation type")
				}
			}
		}
	}()
}

func GetBatch() []BatchOp {
	return batch
}

// KeyValueStoreDB
func (rdb *BatchedRedisDB) Get(k string) (string, error) {
	op := GetOp{K: k}
	ch := make(chan string)
	AppendToBatch(rdb, op, ch)
	result := <-ch

	return result, nil
	// result, err := rdb.rc.Get(context.Background(), k).Result()
	// if err == redis.Nil {
	// 	return "0", nil
	// } else if err != nil {
	// 	log.Printf("Error fetching %s from Redis: %v", k, err)
	// 	return "", err
	// }
	// return result, nil
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
