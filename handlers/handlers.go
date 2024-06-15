package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var (
	posts  = make(map[int]Post)
	nextID = 1
	mu     sync.Mutex
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	post.ID = nextID
	nextID++
	posts[post.ID] = post
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	post, exists := posts[id]
	mu.Unlock()
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	post, exists := posts[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.Author = updatedPost.Author
	posts[id] = post
	mu.Unlock()

	json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	if _, exists := posts[id]; !exists {
		mu.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	delete(posts, id)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func ListPosts(w http.ResponseWriter, r *http.Request) {

	var result []Post

	mu.Lock()	
	for _, post := range posts {
		result = append(result, post)
	}
	mu.Unlock()
	
	json.NewEncoder(w).Encode(result)
}
