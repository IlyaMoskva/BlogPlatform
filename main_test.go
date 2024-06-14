package main

import (
	"blogplatform/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

var (
	posts  = make(map[int]handlers.Post)
	nextID = 1
	mu     sync.Mutex
)

func resetPosts() {
	mu.Lock()
	defer mu.Unlock()
	posts = make(map[int]handlers.Post)
	nextID = 1
}

func TestCreatePost(t *testing.T) {
	resetPosts()

	payload := `{"title": "Test Post", "content": "This is a test post", "author": "Tester"}`
	req, err := http.NewRequest("POST", "/post", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreatePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var got handlers.Post
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("unable to decode response: %v", err)
	}

	expected := handlers.Post{ID: 1, Title: "Test Post", Content: "This is a test post", Author: "Tester"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestGetPost(t *testing.T) {
	resetPosts()

	var posts = make(map[int]handlers.Post)
	posts[1] = handlers.Post{ID: 1, Title: "Test Post", Content: "This is a test post", Author: "Tester"}
	req, err := http.NewRequest("GET", "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetPost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var got handlers.Post
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("unable to decode response: %v", err)
	}

	expected := handlers.Post{ID: 1, Title: "Test Post", Content: "This is a test post", Author: "Tester"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}
