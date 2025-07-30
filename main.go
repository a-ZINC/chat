package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	manager := NewManager()
	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Open("./frontend/index.html"); os.IsNotExist(err) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./frontend/index.html")
	}))
	http.HandleFunc("/room", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Open("./frontend/room.html"); os.IsNotExist(err) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./frontend/room.html")
	}))
	http.HandleFunc("/ws", manager.ServerWs)

	
	fmt.Printf("Server started at http://localhost:3000\n")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
