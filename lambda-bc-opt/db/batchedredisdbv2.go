package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BatchedRedisDBV2 struct {
	batchserviceIP string
}

func (rdb *BatchedRedisDBV2) Get(k string) (string, error) {
	op := GetOp{K: k}

	jsonData, err := json.Marshal(op)
	if err != nil {
		return "", fmt.Errorf("error serializing GetOp: %v", err)
	}

	url := fmt.Sprintf("http://%s/get", rdb.batchserviceIP)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

func ConsBatchedRedisDBV2(batchserviceIP string) *BatchedRedisDBV2 {
	return &BatchedRedisDBV2{
		batchserviceIP: batchserviceIP,
	}
}
