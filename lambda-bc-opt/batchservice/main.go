package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lambda-bc-opt/db"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type BatchedRedisDB struct {
	rc *redis.Client
}

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

type BatchResponse struct {
	redisResponse redis.Cmder
	batchOp       BatchOp
}

var batch chan BatchOp
var batchSize = 10
var loopInterval = 1000 * time.Millisecond

var mu sync.Mutex

var lastExec time.Time

func execPipeline(ctx context.Context, pipe redis.Pipeliner) error {
	for retries := 0; retries < 3; retries++ {
		_, err := pipe.Exec(ctx)
		if err == nil {
			// Pipeline executed successfully
			return nil
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("Pipeline operation timed out. Retrying... attempt %d", retries+1)
		} else if err != nil && err != redis.Nil {
			log.Printf("Pipeline failed with error: %v. Retrying... attempt %d", err, retries+1)
		}

		// Retry with a delay
		time.Sleep(2 * time.Millisecond)
	}
	return fmt.Errorf("pipeline execution failed after 3 retries")
}

func execBatch(rdb *BatchedRedisDB) {
	ctx := context.Background()
	batchResponses := make(chan BatchResponse, batchSize)
	// do the operations in pipleline
	pipe := rdb.rc.Pipeline()
	processedCount := 0
	for i := 0; i < batchSize; i++ {
		select {
		case cur_op := <-batch:
			{
				processedCount++
				switch v := cur_op.op.(type) {
				case GetOp:
					result := pipe.Get(ctx, v.K)
					batchResponses <- BatchResponse{
						redisResponse: result,
						batchOp:       cur_op}
				case SetOp:
					result := pipe.Set(ctx, v.K, v.V, 0)
					batchResponses <- BatchResponse{
						redisResponse: result,
						batchOp:       cur_op}
				default:
					log.Fatalln("Unknown operation type")
				}
			}
		default:
			break
		}
	}
	// executing the pipeline
	if processedCount > 0 {
		log.Printf(">> processedCount %d", processedCount)
		log.Printf(">> batch size %d", len(batch))
		execPipeline(ctx, pipe)
	}
forLoop:
	for {

		select {
		case response := <-batchResponses:
			redisResponse := response.redisResponse
			batchOp := response.batchOp
			switch batchOp.op.(type) {
			case GetOp:
				if getResp, ok := redisResponse.(*redis.StringCmd); ok {
					batchOp.ch <- getResp.Val()
				} else {
					batchOp.ch <- "error"
				}
			case SetOp:
				if setResp, ok := redisResponse.(*redis.StatusCmd); ok {
					batchOp.ch <- setResp.Val()
				} else {
					batchOp.ch <- "error"
				}
			}
		default:
			break forLoop
		}
	}

	if processedCount > 0 {
		curExec := time.Now()
		diff := curExec.Sub(lastExec)
		log.Println("Time difference:", diff)
		lastExec = curExec
	}
}

func appendToBatch(rdb *BatchedRedisDB, op Op, ch chan string) {
	switch op.(type) {
	case GetOp:
		batch <- BatchOp{op, ch}
	case SetOp:
		batch <- BatchOp{op, ch}
	default:
		log.Fatalln("Unknown operation type")
	}
	if len(batch) >= batchSize {
		execBatch(rdb)
		log.Println("exec: BATCH FULL!")
	}
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
		appendToBatch(rdb, getOp, ch)
		result := <-ch

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func main() {
	// DB
	rc := db.InitRedis()
	rdb := BatchedRedisDB{rc: rc}
	batch = make(chan BatchOp, 100*batchSize)

	go func(rdb *BatchedRedisDB) {
		for range time.Tick(loopInterval) {
			if len(batch) > 0 {
				log.Printf("loop: batch size => %d", len(batch))
				log.Println("exec: TL reached!")
				execBatch(rdb)
			}
			// batch = nil
		}
	}(&rdb)

	// API
	http.HandleFunc("/get", getHandler(&rdb))
	fmt.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe("10.10.0.1:8080", nil))
}
