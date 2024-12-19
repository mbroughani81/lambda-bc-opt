package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type BatchedRedisDBV2 struct {
	batchserviceIP string
	batchservicePort string
	buf *bytes.Buffer
	client *http.Client
}


func (rdb *BatchedRedisDBV2) Get(k string) (string, error) {
	start := time.Now()

	op := GetOp{K: k}
	jsonData, err := json.Marshal(op)
	if err != nil {
		return "", fmt.Errorf("error serializing GetOp: %v", err)
	}
	rdb.buf.Reset()
	rdb.buf.Write(jsonData)
	slog.Debug(fmt.Sprintf("buf: %v\n", rdb.buf))

	url := fmt.Sprintf("http://%s:%s/get", rdb.batchserviceIP, rdb.batchservicePort)
	// resp, err := http.Post(url, "application/json", rdb.buf) //use different client for each worker!
	req, err := http.NewRequest(http.MethodPost, url, rdb.buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := rdb.client.Do(req)
	// resp := &http.Response{
	//	Status:     "200 OK",
	//	StatusCode: http.StatusOK,
	//	Header:     make(http.Header),
	//	Body:       io.NopCloser(bytes.NewBufferString("Hello, this is a mock response!")),
	// }
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	end := time.Now()
	slog.Debug(fmt.Sprintf("BatchedRedisDBV2 Get => %v", end.Sub(start)))

	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %s", body)
	}

	return string(body), nil
}

func (rdb *BatchedRedisDBV2) Set(k string, v string) error {
	return errors.New("Not implemented!")
}

func ConsBatchedRedisDBV2(batchserviceIP string, batchservicePort string) *BatchedRedisDBV2 {
	buf := bytes.NewBuffer(make([]byte, 1000))
	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: transport,
	}
	return &BatchedRedisDBV2{
		batchserviceIP: batchserviceIP,
		batchservicePort: batchservicePort,
		buf: buf,
		client: client,
	}
}
