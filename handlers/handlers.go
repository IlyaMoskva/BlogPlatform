package handlers

import (
	"blogplatform/structs"
	"blogplatform/validation"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

// CreatePost creates a new post
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
	Store.PostsList = append(Store.PostsList, post)
	Store.Mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetPost retrieves a post by ID
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
	// Increment the view count
	post.Views++
	Store.Posts[id] = post

	// Increment Global statistics
	structs.GlobalStats.PostViews[id]++
	structs.GlobalStats.AuthorViews[post.Author]++

	json.NewEncoder(w).Encode(post)
}

// UpdatePost updates a post by ID
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

// DeletePost deletes a post by ID
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

// ListPosts lists all posts with pagination
// @Summary List all posts with pagination
// @Description List all blog posts with pagination
// @Tags Post Collection API
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} structs.Post
// @Router /posts [get]
func ListPosts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	Store.Mutex.Lock()
	defer Store.Mutex.Unlock()

	var result []structs.Post
	start := (page - 1) * size
	end := start + size

	if start >= len(Store.PostsList) {
		json.NewEncoder(w).Encode(result)
		return
	}

	if end > len(Store.PostsList) {
		end = len(Store.PostsList)
	}

	result = Store.PostsList[start:end]
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

// SearchPosts searches posts by title or author
// @Summary Search posts by title, content, or author
// @Description Search posts by title, content, or author
// @Tags Post Collection API
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} structs.Post
// @Router /posts/search [get]
func SearchPosts(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if err := validation.ValidateQuery(query); err != nil {
		http.Error(w, err.Error(), err.(validation.HttpError).Code)
		return
	}

	Store.Mutex.Lock()
	defer Store.Mutex.Unlock()

	var result []structs.Post
	for _, post := range Store.Posts {
		if strings.Contains(post.Title, query) || strings.Contains(post.Content, query) || strings.Contains(post.Author, query) {
			post.SearchAppearances++
			Store.Posts[post.ID] = post
			result = append(result, post)
		}
	}
	// Update global statistics
	structs.GlobalStats.TotalSearchCount++

	json.NewEncoder(w).Encode(result)
}

// GetReports 	 Provides statistics of views, top authors and posts
// @Summary      Get report statistics
// @Description  Retrieve various report statistics such as number of views for posts, number of search appearances, top authors, and top posts.
// @Tags         Reports API
// @Produce      json
// @Success      200  {object}  structs.ReportResponse
// @Failure      400  {object}  structs.ErrorResponse
// @Failure      500  {object}  structs.ErrorResponse
// @Router       /reports [get]
func GetReports(w http.ResponseWriter, r *http.Request) {
	reports := structs.ReportResponse{
		TotalViews:        getTotalViews(),
		SearchAppearances: getSearchAppearances(),
		TopAuthors:        getTopAuthors(),
		TopPosts:          getTopPosts(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func getTotalViews() int {
	return sumMapValues(structs.GlobalStats.PostViews)
}

func getSearchAppearances() int {
	return structs.GlobalStats.TotalSearchCount
}

func getTopAuthors() []structs.Author {
	return getTopItems(structs.GlobalStats.AuthorViews)
}

func getTopPosts() []structs.Post {
	topPosts := make([]structs.Post, 0, len(structs.GlobalStats.PostViews))
	for id, views := range structs.GlobalStats.PostViews {
		post, exists := Store.Posts[id]
		if exists {
			post.Views = views
			topPosts = append(topPosts, post)
		}
	}
	sort.Slice(topPosts, func(i, j int) bool {
		return topPosts[i].Views > topPosts[j].Views
	})
	return topPosts
}

func sumMapValues(m map[int]int) int {
	total := 0
	for _, value := range m {
		total += value
	}
	return total
}

func getTopItems(m map[string]int) []structs.Author {
	items := make([]structs.Author, 0, len(m))
	for name, count := range m {
		items = append(items, structs.Author{Name: name, Views: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Views == items[j].Views {
			return items[i].Name < items[j].Name
		}
		return items[i].Views > items[j].Views
	})
	return items
}
