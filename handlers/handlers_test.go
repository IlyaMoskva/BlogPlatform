package handlers

import (
	"blogplatform/structs"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createDummyPost(payload string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	return rr
}

func resetPosts() {
	Store.Mutex.Lock()
	defer Store.Mutex.Unlock()
	Store.Posts = make(map[int]structs.Post)
	Store.PostsList = []structs.Post{}
}

func resetStats() {
	structs.GlobalStats = structs.Stats{
		PostViews:         make(map[int]int),
		SearchAppearances: make(map[int]int),
		AuthorViews:       make(map[string]int),
		TotalSearchCount:  0,
	}
}

func TestCreatePost(t *testing.T) {
	payload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	rr := createDummyPost(payload, t)

	var post structs.Post
	if err := json.NewDecoder(rr.Body).Decode(&post); err != nil {
		t.Fatal(err)
	}

	if post.Title != "Test Post" || post.Content != "This is a test post" || post.Author != "Tester" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestGetPost(t *testing.T) {
	// Simulate creating a post first
	payload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createDummyPost(payload, t)

	// Now test retrieving the post
	getReq, err := http.NewRequest(http.MethodGet, "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(GetPost)
	getHandler.ServeHTTP(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var post structs.Post
	if err := json.NewDecoder(getRR.Body).Decode(&post); err != nil {
		t.Fatal(err)
	}

	expected := structs.Post{
		ID:                1,
		Title:             "Test Post",
		Content:           "This is a test post",
		Author:            "Tester",
		Views:             1,
		SearchAppearances: 0,
	}
	if post != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			post, expected)
	}
}

func TestUpdatePost(t *testing.T) {
	resetPosts()
	// Simulate creating a post first
	payload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createDummyPost(payload, t)

	// Now test updating the post
	updatePayload := `{"title":"Updated Post","content":"This is an updated test post","author":"Tester"}`
	updateReq, err := http.NewRequest(http.MethodPut, "/post?id=1", bytes.NewBuffer([]byte(updatePayload)))
	if err != nil {
		t.Fatal(err)
	}

	updateRR := httptest.NewRecorder()
	updateHandler := http.HandlerFunc(UpdatePost)
	updateHandler.ServeHTTP(updateRR, updateReq)

	if status := updateRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedPost structs.Post
	if err := json.NewDecoder(updateRR.Body).Decode(&updatedPost); err != nil {
		t.Fatal(err)
	}

	expected := structs.Post{
		ID:      1,
		Title:   "Updated Post",
		Content: "This is an updated test post",
		Author:  "Tester",
	}
	if updatedPost != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			updatedPost, expected)
	}
}

func TestDeletePost(t *testing.T) {
	resetPosts()
	// Simulate creating a post first
	payload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createDummyPost(payload, t)

	// Now test deleting the post
	deleteReq, err := http.NewRequest(http.MethodDelete, "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	deleteRR := httptest.NewRecorder()
	deleteHandler := http.HandlerFunc(DeletePost)
	deleteHandler.ServeHTTP(deleteRR, deleteReq)

	if status := deleteRR.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestListPostsWithPagination(t *testing.T) {
	resetPosts()
	// Simulate creating two posts
	payload1 := `{"id":1,"title":"Test Post1","content":"This is a test post","author":"Tester1"}`
	payload2 := `{"id":2,"title":"Test Post2","content":"This is a test post","author":"Tester2"}`
	createDummyPost(payload1, t)
	createDummyPost(payload2, t)

	if len(Store.Posts) != 2 {
		t.Fatalf("expected 2 posts in the store, got %d", len(Store.Posts))
	}

	req, err := http.NewRequest("GET", "/posts?page=1&size=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListPosts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var posts []structs.Post
	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestImportPostsFromFile(t *testing.T) {
	resetPosts()
	// Prepare the file
	path := filepath.Join("../testfiles", "test_posts.json")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Prepare the multipart form file
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = file.Seek(0, 0) // Reset file pointer to the beginning
	if err != nil {
		t.Fatalf("Failed to seek file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file content: %v", err)
	}
	writer.Close()

	// Create the request
	req := httptest.NewRequest("POST", "/admin/import", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create the response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	ImportPostsFromFile(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Posts imported successfully"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Verify the in-memory store
	if len(Store.Posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(Store.Posts))
	}

	// Check if the posts were correctly imported
	post1, exists := Store.Posts[1]
	if !exists || post1.Title != "Title 1" || post1.Content != "Content 1" || post1.Author != "Author 1" {
		t.Errorf("Post 1 not imported correctly: %+v", post1)
	}

	post2, exists := Store.Posts[2]
	if !exists || post2.Title != "Title 2" || post2.Content != "Content 2" || post2.Author != "Author 2" {
		t.Errorf("Post 2 not imported correctly: %+v", post2)
	}
}

func TestSearchPosts(t *testing.T) {
	resetPosts()
	// Simulate creating multiple posts
	payload1 := `{"title":"First Post","content":"This is the first test post","author":"Author1"}`
	payload2 := `{"title":"Second Post","content":"This is the second test post","author":"Author2"}`
	payload3 := `{"title":"Another Post","content":"Another post for testing","author":"Author3"}`
	createDummyPost(payload1, t)
	createDummyPost(payload2, t)
	createDummyPost(payload3, t)

	if len(Store.Posts) != 3 {
		t.Fatalf("expected 3 posts in the store, got %d", len(Store.Posts))
	}

	// Define test cases
	tests := []struct {
		query    string
		expected int
	}{
		{"First", 1},
		{"Post", 3},
		{"test", 3},
		{"Author1", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/search?q="+tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SearchPosts)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			var posts []structs.Post
			if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
				t.Fatal(err)
			}

			if len(posts) != tt.expected {
				t.Errorf("expected %d posts, got %d for query %v", tt.expected, len(posts), tt.query)
			}
		})
	}
}

func TestGetReports(t *testing.T) {
	resetPosts()
	resetStats()

	// Simulate creating posts
	payload1 := `{"title":"Test Post 1","content":"This is the first test post","author":"Tester1"}`
	payload2 := `{"title":"Test Post 2","content":"This is the second test post","author":"Tester2"}`
	createDummyPost(payload1, t)
	createDummyPost(payload2, t)

	// Simulate viewing posts
	getReq1, err := http.NewRequest(http.MethodGet, "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	getRR1 := httptest.NewRecorder()
	getHandler := http.HandlerFunc(GetPost)
	getHandler.ServeHTTP(getRR1, getReq1)

	getReq2, err := http.NewRequest(http.MethodGet, "/post?id=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	getRR2 := httptest.NewRecorder()
	getHandler.ServeHTTP(getRR2, getReq2)

	// Simulate searching for posts
	searchReq, err := http.NewRequest(http.MethodGet, "/posts?q=Test", nil)
	if err != nil {
		t.Fatal(err)
	}
	searchRR := httptest.NewRecorder()
	searchHandler := http.HandlerFunc(SearchPosts)
	searchHandler.ServeHTTP(searchRR, searchReq)

	// Now test getting the reports
	reportsReq, err := http.NewRequest(http.MethodGet, "/reports", nil)
	if err != nil {
		t.Fatal(err)
	}

	reportsRR := httptest.NewRecorder()
	reportsHandler := http.HandlerFunc(GetReports)
	reportsHandler.ServeHTTP(reportsRR, reportsReq)

	if status := reportsRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var reports structs.ReportResponse
	if err := json.NewDecoder(reportsRR.Body).Decode(&reports); err != nil {
		t.Fatal(err)
	}

	expectedTotalViews := 2
	if reports.TotalViews != expectedTotalViews {
		t.Errorf("handler returned unexpected total views: got %v want %v", reports.TotalViews, expectedTotalViews)
	}

	expectedSearchAppearances := 1
	if reports.SearchAppearances != expectedSearchAppearances {
		t.Errorf("handler returned unexpected search appearances: got %v want %v", reports.SearchAppearances, expectedSearchAppearances)
	}

	expectedTopAuthors := []structs.Author{
		{Name: "Tester1", Views: 1},
		{Name: "Tester2", Views: 1},
	}
	if !compareAuthors(reports.TopAuthors, expectedTopAuthors) {
		t.Errorf("handler returned unexpected top authors: got %v want %v", reports.TopAuthors, expectedTopAuthors)
	}

	expectedTopPosts := []structs.Post{
		{ID: 1, Title: "Test Post 1", Content: "This is the first test post", Author: "Tester1", Views: 1},
		{ID: 2, Title: "Test Post 2", Content: "This is the second test post", Author: "Tester2", Views: 1},
	}
	if !comparePosts(reports.TopPosts, expectedTopPosts) {
		t.Errorf("handler returned unexpected top posts: got %v want %v", reports.TopPosts, expectedTopPosts)
	}
}

func compareAuthors(a, b []structs.Author) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Views != b[i].Views {
			return false
		}
	}
	return true
}

func comparePosts(a, b []structs.Post) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ID != b[i].ID || a[i].Title != b[i].Title || a[i].Content != b[i].Content || a[i].Author != b[i].Author || a[i].Views != b[i].Views {
			return false
		}
	}
	return true
}
