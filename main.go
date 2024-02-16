package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!!!!")
	})

	fmt.Println("Server is running on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
