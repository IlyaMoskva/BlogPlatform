package handlers

import (
	_ "blogplatform/docs"
	"blogplatform/validation"
	"encoding/json"
	"net/http"
	"sync"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var (
	Posts  = make(map[int]Post)
	nextID = 1
	mu     sync.Mutex
)

func extractAndValidateID(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	return validation.ValidateID(idStr)
}

// @title Blog API
// @version 1.0
// @description This is a sample server for a blog.
// @host localhost:8443
// @BasePath /

// createPost creates a new post
// @Summary Create a new post
// @Description Create a new blog post
// @Tags Post API
// @Accept json
// @Produce json
// @Param post body Post true "Post content"
// @Success 201 {object} Post
// @Router /post [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	post.ID = nextID
	nextID++
	Posts[post.ID] = post
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// getPost retrieves a post by ID
// @Summary Get a post by ID
// @Description Get a post by ID
// @Tags Post API
// @Accept json
// @Produce json
// @Param id query int true "Post ID"
// @Success 200 {object} Post
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Router /post [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := extractAndValidateID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	post, exists := Posts[id]
	mu.Unlock()
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// updatePost updates a post by ID
// @Summary Update a post by ID
// @Description Update a post by ID
// @Tags Post API
// @Accept json
// @Produce json
// @Param id query int true "Post ID"
// @Param post body Post true "Post content"
// @Success 200 {object} Post
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Router /post [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := extractAndValidateID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	post, exists := Posts[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.Author = updatedPost.Author
	Posts[id] = post
	mu.Unlock()

	json.NewEncoder(w).Encode(post)
}

// deletePost deletes a post by ID
// @Summary Delete a post by ID
// @Description Delete a post by ID
// @Tags Post API
// @Accept json
// @Produce json
// @Param id query int true "Post ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Router /post [delete]
func DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := extractAndValidateID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	if _, exists := Posts[id]; !exists {
		mu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	delete(Posts, id)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// listPosts lists all posts
// @Summary List all posts
// @Description List all blog posts
// @Tags Post API
// @Accept json
// @Produce json
// @Success 200 {array} Post
// @Router /posts [get]
func ListPosts(w http.ResponseWriter, r *http.Request) {
	var result []Post

	mu.Lock()
	for _, post := range Posts {
		result = append(result, post)
	}
	mu.Unlock()

	json.NewEncoder(w).Encode(result)
}
