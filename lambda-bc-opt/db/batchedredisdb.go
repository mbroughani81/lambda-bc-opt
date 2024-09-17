package db

import (
	"context"
	"log"
	"sync"
	"time"

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
var batchSize = 100
var loopInterval = 2 * time.Second

var mu sync.Mutex

func ExecBatch(rdb *BatchedRedisDB) {
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

func AppendToBatch(rdb *BatchedRedisDB, op Op, ch chan string) {
	go func() {
		// critical section! only one coroutine here
		mu.Lock()
		defer mu.Unlock()

		switch v := op.(type) {
		case GetOp:
			log.Println("Appending GetOp:", v.K)
			batch = append(batch, BatchOp{op, ch}) // Append the GetOp to the batch
		case SetOp:
			log.Println("Appending SetOp:", v.K, v.V)
			batch = append(batch, BatchOp{op, ch}) // Append the SetOp to the batch
		default:
			log.Fatalln("Unknown operation type")
		}
		if len(batch) >= batchSize {
			ExecBatch(rdb)
			// delete the batch
			batch = nil
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
}
func (rdb *BatchedRedisDB) Set(k string, v string) error {
	op := SetOp{K: k, V: v}
	ch := make(chan string)
	AppendToBatch(rdb, op, ch)
	<-ch
	return nil
}

func ConsBatchedRedisDB() *BatchedRedisDB {
	rc := initRedis()
	rdb := BatchedRedisDB{rc: rc}

	go func(rdb *BatchedRedisDB) {
		for now := range time.Tick(loopInterval) {
			mu.Lock()
			ExecBatch(rdb)
			batch = nil
			log.Printf("Executing batch => %s", now)
			mu.Unlock()
		}
	}(&rdb)

	return &rdb
}
