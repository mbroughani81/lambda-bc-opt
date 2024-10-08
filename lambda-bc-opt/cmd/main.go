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
	rdb := db.ConsRedisDB("10.10.0.1", "6379")

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
	rdb := db.ConsRedisDB("10.10.0.1", "6379")

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

const numWorkers = 5

// Worker function that simulates a blocking operation (e.g., network request)
func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		// Simulate blocking with a network request
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Worker %d: failed to fetch %s: %v\n", id, url, err)
			continue
		}
		fmt.Printf("Worker %d: fetched %s with status %s\n", id, url, resp.Status)
		resp.Body.Close()
		time.Sleep(1 * time.Second) // Simulate some additional blocking time
	}
}

func fn4() {
	var wg sync.WaitGroup
	jobs := make(chan string, 10)

	// Start a pool of workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Enqueue jobs
	urls := []string{
		"https://golang.org",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://golang.org",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://golang.org",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
	}

	for _, url := range urls {
		jobs <- url
	}

	close(jobs)
	wg.Wait() // Wait for all workers to finish
	fmt.Println("All workers done")
}

func main() {
	log.Println("START")
	fn4()
}
