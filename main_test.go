package main

import (
	"blogplatform/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePost(t *testing.T) {
	payload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`

	req, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreatePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var post handlers.Post
	if err := json.NewDecoder(rr.Body).Decode(&post); err != nil {
		t.Fatal(err)
	}

	if post.Title != "Test Post" || post.Content != "This is a test post" || post.Author != "Tester" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestGetPost(t *testing.T) {
	// Simulate creating a post first
	createPayload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createReq, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(createPayload)))
	if err != nil {
		t.Fatal(err)
	}
	createRR := httptest.NewRecorder()
	createHandler := http.HandlerFunc(handlers.CreatePost)
	createHandler.ServeHTTP(createRR, createReq)

	// Now test retrieving the post
	getReq, err := http.NewRequest(http.MethodGet, "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(handlers.GetPost)
	getHandler.ServeHTTP(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var post handlers.Post
	if err := json.NewDecoder(getRR.Body).Decode(&post); err != nil {
		t.Fatal(err)
	}

	expected := handlers.Post{
		ID:      1,
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Tester",
	}
	if post != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			post, expected)
	}
}

func TestUpdatePost(t *testing.T) {
	// Simulate creating a post first
	createPayload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createReq, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(createPayload)))
	if err != nil {
		t.Fatal(err)
	}
	createRR := httptest.NewRecorder()
	createHandler := http.HandlerFunc(handlers.CreatePost)
	createHandler.ServeHTTP(createRR, createReq)

	// Now test updating the post
	updatePayload := `{"title":"Updated Post","content":"This is an updated test post","author":"Tester"}`
	updateReq, err := http.NewRequest(http.MethodPut, "/post?id=1", bytes.NewBuffer([]byte(updatePayload)))
	if err != nil {
		t.Fatal(err)
	}

	updateRR := httptest.NewRecorder()
	updateHandler := http.HandlerFunc(handlers.UpdatePost)
	updateHandler.ServeHTTP(updateRR, updateReq)

	if status := updateRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedPost handlers.Post
	if err := json.NewDecoder(updateRR.Body).Decode(&updatedPost); err != nil {
		t.Fatal(err)
	}

	expected := handlers.Post{
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
	// Simulate creating a post first
	createPayload := `{"title":"Test Post","content":"This is a test post","author":"Tester"}`
	createReq, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte(createPayload)))
	if err != nil {
		t.Fatal(err)
	}
	createRR := httptest.NewRecorder()
	createHandler := http.HandlerFunc(handlers.CreatePost)
	createHandler.ServeHTTP(createRR, createReq)

	// Now test deleting the post
	deleteReq, err := http.NewRequest(http.MethodDelete, "/post?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	deleteRR := httptest.NewRecorder()
	deleteHandler := http.HandlerFunc(handlers.DeletePost)
	deleteHandler.ServeHTTP(deleteRR, deleteReq)

	if status := deleteRR.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}
