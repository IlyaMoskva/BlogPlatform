package main

import (
	_ "blogplatform/docs"
	"blogplatform/handlers"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

var version float64 = 0.5

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

	// Swagger UI endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server started at https://localhost:8443")

	// Serve HTTPS with TLS
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
