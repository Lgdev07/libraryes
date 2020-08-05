package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lgdev07/libraryes/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRouteCreateBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedOneAuthor()
	if err != nil {
		log.Fatal(err)
	}

	inputJSON := `{
		"name": "Book1",
		"edition": 1,
		"publication_year": 2020,
		"authors": [1]
	}`

	req, err := http.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreateBook)
	handler.ServeHTTP(rr, req)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusCreated)
	assert.Equal(t, responseInterface["status"], "success")

	bookResponse := responseInterface["book"].(map[string]interface{})

	assert.Equal(t, bookResponse["name"], "Book1")
	assert.Equal(t, bookResponse["ID"], float64(1))
	assert.Equal(t, bookResponse["edition"], float64(1))
	assert.Equal(t, bookResponse["publication_year"], float64(2020))

	authorResponse := bookResponse["authors"].([]interface{})[0].(map[string]interface{})

	assert.Equal(t, authorResponse["ID"], float64(1))
	assert.Equal(t, authorResponse["name"], "João Doe")

}

func TestRouteDeleteBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	book, err := seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, book.ID, uint(1))

	req, err := http.NewRequest(http.MethodDelete, "/books/"+"1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", app.DeleteBook)
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	var params map[string]string

	books, err := models.GetAllBooks(app.DB, params)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*books), 0)
}

func TestRouteUpdateBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	book, err := seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, book.ID, uint(1))
	assert.Equal(t, book.Name, "Book 1")

	err = seedOneAuthor()
	if err != nil {
		log.Fatal(err)
	}

	inputJSON := `{
		"name": "Book 2",
		"edition": 1,
		"publication_year": 2020,
		"authors": [1]
	}`

	req, err := http.NewRequest(http.MethodPut, "/books/"+"1", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", app.UpdateBook)
	router.ServeHTTP(rr, req)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, responseInterface["status"], "success")

	bookResponse := responseInterface["book"].(map[string]interface{})

	assert.Equal(t, bookResponse["name"], "Book 2")
	assert.Equal(t, bookResponse["ID"], float64(1))
	assert.Equal(t, bookResponse["edition"], float64(1))
	assert.Equal(t, bookResponse["publication_year"], float64(2020))

	authorResponse := bookResponse["authors"].([]interface{})[0].(map[string]interface{})

	assert.Equal(t, authorResponse["ID"], float64(1))
	assert.Equal(t, authorResponse["name"], "João Doe")

}

func TestRouteListBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	book, err := seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, book.ID, uint(1))
	assert.Equal(t, book.Name, "Book 1")

	// When we don't pass a query param it should return all books
	req, err := http.NewRequest(http.MethodGet, "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ListBooks)
	handler.ServeHTTP(rr, req)

	var responseInterface []map[string]interface{}

	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(responseInterface), 1)
	assert.Equal(t, responseInterface[0]["name"], "Book 1")

	// When the query param is incorrect it should return nothing
	req, err = http.NewRequest(http.MethodGet, "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("name", "Wrong Name")

	req.URL.RawQuery = q.Encode()

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(app.ListBooks)
	handler.ServeHTTP(rr, req)

	var response2Interface []map[string]interface{}

	err = json.Unmarshal([]byte(rr.Body.String()), &response2Interface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(response2Interface), 0)

}
