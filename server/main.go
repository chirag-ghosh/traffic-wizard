package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = fmt.Sprintf("Hello from Server: %s", os.Getenv("id"))
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error in JSON marshal: %s", err)
	}
	w.Write(jsonResp)
	return
}

func heartbeatEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return
}

func main() {
	http.HandleFunc("/home", homeEndpoint)
	http.HandleFunc("/heartbeat", heartbeatEndpoint)

	fmt.Println("Starting server on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		panic(err)
	}
}
