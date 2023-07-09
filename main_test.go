package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetBooks(t *testing.T) {
	Books = append(Books, Book{ID: "1", Title: "Book One", Author: "Author One"})
	Books = append(Books, Book{ID: "2", Title: "Book Two", Author: "Author Two"})
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"1","title":"Book One","author":"Author One"},{"id":"2","title":"Book Two","author":"Author Two"}]`
	if strings.ReplaceAll(rr.Body.String(), "\n", "") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
