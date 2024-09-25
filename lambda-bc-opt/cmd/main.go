package main

import (
	"fmt"
	"lambda-bc-opt/db"
	"log"
	"strconv"
	"sync"
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

func main() {
	log.Println("START")
	fn2()
}
