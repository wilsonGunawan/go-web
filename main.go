package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Book struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}

var Books []Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range Books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	Books = append(Books, book)
	json.NewEncoder(w).Encode(book)
}

func Sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//loop 10 times
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		fmt.Fprintf(w, "data: %s\n\n", "hello")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func main() {
	router := mux.NewRouter()

	Books = append(Books, Book{ID: "1", Title: "Book One", Author: "Author One"})
	Books = append(Books, Book{ID: "2", Title: "Book Two", Author: "Author Two"})

	router.HandleFunc("/Books", GetBooks).Methods("GET")
	router.HandleFunc("/Books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/Books", CreateBook).Methods("POST")
	router.HandleFunc("/sse", Sse).Methods("GET")

	//db := pg.Connect(&pg.Options{User: "postgres"})
	//defer db.Close()

	bytes, err := bcrypt.GenerateFromPassword([]byte("testing12345"), 14)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	err = bcrypt.CompareHashAndPassword(bytes, []byte("testing12345"))
	fmt.Println(err)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
