package main

import (
	"blogplatform/handlers"
	"fmt"
	"net/http"
)

var version float64 = 0.3

func main() {
	fmt.Printf("Blog platform API. v %v\n", version)

	http.HandleFunc("/posts", handlers.ListPosts)
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreatePost(w, r)
		case http.MethodGet:
			handlers.GetPost(w, r)
		case http.MethodPut:
			handlers.UpdatePost(w, r)
		case http.MethodDelete:
			handlers.DeletePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
