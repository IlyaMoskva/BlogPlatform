package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportPostsFromFile(t *testing.T) {
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
