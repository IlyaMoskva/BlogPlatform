package handlers

import (
	"blogplatform/structs"
	"blogplatform/validation"
	"encoding/json"
	"net/http"
	"sync"
)

type InMemoryStore struct {
	Posts     map[int]structs.Post
	PostsList []structs.Post
	Mutex     sync.Mutex
}

var Store = InMemoryStore{
	Posts:     make(map[int]structs.Post),
	PostsList: []structs.Post{},
}

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
// @Param post body structs.Post true "Post content"
// @Success 201 {object} structs.Post
// @Router /post [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post structs.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Store.Mutex.Lock()
	post.ID = len(Store.PostsList) + 1
	Store.Posts[post.ID] = post
	Store.Mutex.Unlock()

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
// @Success 200 {object} structs.Post
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Router /post [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := extractAndValidateID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Store.Mutex.Lock()
	post, exists := Store.Posts[id]
	Store.Mutex.Unlock()
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
// @Param post body structs.Post true "Post content"
// @Success 200 {object} structs.Post
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Router /post [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := extractAndValidateID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedPost structs.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Store.Mutex.Lock()
	defer Store.Mutex.Unlock()

	post, exists := Store.Posts[id]
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.Author = updatedPost.Author
	Store.Posts[id] = post

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

	Store.Mutex.Lock()
	if _, exists := Store.Posts[id]; !exists {
		Store.Mutex.Unlock()
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	delete(Store.Posts, id)
	Store.Mutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// listPosts lists all posts
// @Summary List all posts
// @Description List all blog posts
// @Tags Post API
// @Accept json
// @Produce json
// @Success 200 {array} structs.Post
// @Router /posts [get]
func ListPosts(w http.ResponseWriter, r *http.Request) {
	var result []structs.Post

	Store.Mutex.Lock()
	for _, post := range Store.Posts {
		result = append(result, post)
	}
	Store.Mutex.Unlock()

	json.NewEncoder(w).Encode(result)
}

// ImportPostsFromFile godoc
// @Summary Import posts from a JSON file
// @Description Upload and import posts from a JSON file
// @Tags Admin API
// @Accept multipart/form-data
// @Produce plain
// @Param file formData file true "JSON file with posts"
// @Success 200 {string} string "Posts imported successfully"
// @Failure 400 {string} string "Error retrieving the file"
// @Failure 500 {string} string "Error decoding JSON file"
// @Router /admin/import [post]
func ImportPostsFromFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var data struct {
		Posts []structs.Post `json:"posts"`
	}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding JSON file", http.StatusBadRequest)
		return
	}

	Store.Mutex.Lock()
	defer Store.Mutex.Unlock()

	for _, post := range data.Posts {
		post.ID = len(Store.Posts) + 1
		Store.Posts[post.ID] = post
		Store.PostsList = append(Store.PostsList, post)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Posts imported successfully"))
}
