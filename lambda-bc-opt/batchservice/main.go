package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"lambda-bc-opt/db"
	"log"
	"net/http"
	"sync"
	"time"
)

type BatchedRedisDB struct {
	rc *redis.Client
}

// todo: use generic operations
type Op interface{}

type GetOp struct {
	K string `json:"k"`
}

type SetOp struct {
	K string `json:"k"`
	V string `json:"v"`
}

type BatchOp struct {
	op Op
	ch chan string
}

var batch []BatchOp
var batchSize = 100
var loopInterval = 2000 * time.Millisecond

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
	log.Printf("size of batch => %d", len(batch))
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

// ------------------------------------------------------------ //
func getHandler(rdb *BatchedRedisDB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		// Deserialize the JSON to GetOp struct
		var getOp GetOp
		err = json.Unmarshal(body, &getOp)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received key: %s\n", getOp.K)
		// creating a BatchOp
		ch := make(chan string)
		AppendToBatch(rdb, getOp, ch)
		result := <-ch

		response := fmt.Sprintf("Value for key '%s'", result)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

func main() {
	// DB
	rc := db.InitRedis()
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
	// API
	http.HandleFunc("/get", getHandler(&rdb))
	fmt.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe("10.10.0.1:8080", nil))
}
