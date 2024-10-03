package main

import (
	"fmt"
	"lambda-bc-opt/db"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"net/http"
)

var n int = 200000

func fn1() {
	rdb := db.ConsRedisDB()

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			rdb.Get("cnt")
			wg.Done()

		}()
	}
	wg.Wait()

	cntVal := "1"
	cnt, _ := strconv.Atoi(cntVal)
	cnt++
	greeting := fmt.Sprintf("Hello! cnt is %d.\n", cnt)
	log.Printf("greeting => %s", greeting)

}

func fn2() {
	rdb := db.ConsRedisDB()

	for i := 0; i < n; i++ {
		rdb.Get("cnt")
	}

	cntVal := "1"

	cnt, _ := strconv.Atoi(cntVal)
	cnt++
	greeting := fmt.Sprintf("Hello! cnt is %d.\n", cnt)
	log.Printf("greeting => %s", greeting)
}

func handler(http.ResponseWriter, *http.Request) {
	time.Sleep(3 * time.Second)
	log.Println("sleep done!")
}

func fn3() {
	runtime.GOMAXPROCS(2)
	http.HandleFunc("/get", handler)
	fmt.Println("Server listening on localhost:8080")
	log.Fatal(http.ListenAndServe("10.10.0.1:8080", nil))
}

func main() {
	log.Println("START")
	fn3()
}
