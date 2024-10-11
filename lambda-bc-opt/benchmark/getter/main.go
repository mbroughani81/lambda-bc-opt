import (
	// "io"
	"log"
	"net/http"
	"sync"

	"lambda-bc-opt/db"
)

// var rdb db.KeyValueStoreDB = db.ConsMockRedisDB()

// var rdb db.KeyValueStoreDB = db.ConsRedisDB()

var rdb db.KeyValueStoreDB = db.ConsBatchedRedisDB("10.10.0.1", "6379")

var counter int = 0

func getterHandler(w http.ResponseWriter, r *http.Request) {
	// Call the getter function
	log.Printf("counter => %d", counter)
	counter++
	n := 100

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			rdb.Get("cnt")
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("wg done")
	// Send a response back to the client
}

func main() {
	// log.SetOutput(io.Discard)
	http.HandleFunc("/getter", getterHandler)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
