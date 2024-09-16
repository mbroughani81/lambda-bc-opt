package main

import (
	"hello-world/db"
	"log"
	"sync"
)

func main() {
	batchedRedisDB := db.ConsBatchedRedisDB()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {		
		batchedRedisDB.Set("g1", "value1")
		result, _ := batchedRedisDB.Get("g1")
		log.Println("#1 value1 => %s", result)
		wg.Done()
	}()

	go func() {		
		result, _ := batchedRedisDB.Get("g2")
		log.Println("#2 => %s", result)
		batchedRedisDB.Set("g2", "value2")
		result, _ = batchedRedisDB.Get("g2")
		log.Println("#2 value2 => %s", result)
		wg.Done()
	}()

	// go func() {		
	// 	result, _ := batchedRedisDB.Get("g1")
	// 	log.Println("#1 => %s", result)
	// 	wg.Done()
	// }()


	// go func() {		
	// 	result, _ := batchedRedisDB.Get("g2")
	// 	log.Println("#2 => %s", result)
	// 	wg.Done()
	// }()


	// go func() {		
	// 	result, _ := batchedRedisDB.Get("g3")
	// 	log.Println("#3 => %s", result)
	// 	wg.Done()
	// }()

	wg.Wait()
	// xx := []int{1,2,3}
	// xx = append(xx, 4)
	// log.Println("%#v", xx)
	// xx = nil
	// log.Println("%#v", xx)
	// xx = append(xx, 1, 2, 3, 4, 5, 6)
	// log.Println("%#v", xx)
}

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	go func() {
// 		for now := range time.Tick(5 * time.Second) {
// 			fmt.Println(now, "gooz")
// 		}
// 		wg.Done()
// 	}()

// 	wg.Wait()
// }