package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Lgdev07/olist_challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestRouteCreateAuthor(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	inputJSON := `{"name": "Author 1"}`

	req, err := http.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreateAuthor)
	handler.ServeHTTP(rr, req)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusCreated)
	assert.Equal(t, responseInterface["status"], "success")

	authorResponse := responseInterface["author"].(map[string]interface{})

	assert.Equal(t, authorResponse["name"], "Author 1")
	assert.Equal(t, authorResponse["ID"], float64(1))
}

func TestRouteListAuthors(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedOneAuthor()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "/authors", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ListAuthors)
	handler.ServeHTTP(rr, req)

	var responseInterface []map[string]interface{}

	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, responseInterface[0]["name"], "João Doe")
	assert.Equal(t, responseInterface[0]["ID"], float64(1))
}

func TestUploadImage(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	b, w := createMultipartFormData(t, "file", "authors.csv")
	request := httptest.NewRequest(http.MethodPost, "/authors/import", &b)
	request.Header.Add("Content-Type", w.FormDataContentType())

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ImportCsvAuthor)
	handler.ServeHTTP(response, request)

	var responseInterface map[string]interface{}

	err = json.Unmarshal([]byte(response.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, response.Code, http.StatusCreated)
	assert.Equal(t, responseInterface["status"], "success")
	assert.Equal(t, responseInterface["message"], "Authors successfully created")

	var params map[string]string

	authors, _ := models.GetAllAuthors(app.DB, params)
	assert.Equal(t, len(*authors), 6)

}

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}
