package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
)

type GetOp struct {
    K string `json:"k"`
}

func getHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read the body of the request
    body, err := ioutil.ReadAll(r.Body)
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

    // Log the key received (for demonstration)
    fmt.Printf("Received key: %s\n", getOp.K)

    // Dummy response: Return the value corresponding to the key
    // In a real scenario, you'd retrieve the value from a database or store.
    response := fmt.Sprintf("Value for key '%s'", getOp.K)

    // Send the response back to the client
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}

func main() {
    http.HandleFunc("/get", getHandler)

    fmt.Println("Server listening on localhost:8080")
    // Start the HTTP server on port 8080
    log.Fatal(http.ListenAndServe("10.10.0.1:8080", nil))
}
