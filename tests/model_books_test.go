package tests

import (
	"log"
	"testing"

	"github.com/Lgdev07/libraryes/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	var params map[string]string

	book, err := models.GetAllBooks(app.DB, params)
	if err != nil {
		t.Errorf("There was an error when getting a book, err: %v\n", err)
		return
	}

	assert.Equal(t, len(*book), 1)
	assert.Equal(t, (*book)[0].Name, "Book 1")
	assert.Equal(t, (*book)[0].Edition, uint32(1))
	assert.Equal(t, (*book)[0].PublicationYear, uint32(2020))
	assert.Equal(t, len((*book)[0].Authors), 2)
}

func TestGetBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	book, _ := models.ShowBook(app.DB, 1)

	assert.Equal(t, book.Name, "Book 1")
	assert.Equal(t, book.Edition, uint32(1))
	assert.Equal(t, book.PublicationYear, uint32(2020))

	book, _ = models.ShowBook(app.DB, 2)
	assert.Equal(t, book.Name, "")
	assert.Equal(t, book.Edition, uint32(0))
	assert.Equal(t, book.PublicationYear, uint32(0))
}

func TestDeleteBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	book, _ := models.ShowBook(app.DB, 1)

	assert.Equal(t, book.Name, "Book 1")
	assert.Equal(t, book.Edition, uint32(1))
	assert.Equal(t, book.PublicationYear, uint32(2020))

	err = models.DeleteBook(app.DB, int(book.ID))
	if err != nil {
		log.Fatal(err)
	}

	book, _ = models.ShowBook(app.DB, 1)
	assert.Equal(t, book.Name, "")
	assert.Equal(t, book.Edition, uint32(0))
	assert.Equal(t, book.PublicationYear, uint32(0))
}

func TestGetAllBooks(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	params := map[string]string{
		"name": "Wrong",
	}

	books, err := models.GetAllBooks(app.DB, params)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*books), 0)

	params["name"] = "Book 1"

	books, err = models.GetAllBooks(app.DB, params)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*books), 1)
	assert.Equal(t, (*books)[0].Name, "Book 1")
	assert.Equal(t, (*books)[0].Edition, uint32(1))
	assert.Equal(t, (*books)[0].PublicationYear, uint32(2020))
	assert.Equal(t, len((*books)[0].Authors), 2)

	params["name"] = ""
	params["author"] = "João Doe"

	books, err = models.GetAllBooks(app.DB, params)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*books), 1)
	assert.Equal(t, (*books)[0].Name, "Book 1")
	assert.Equal(t, (*books)[0].Edition, uint32(1))
	assert.Equal(t, (*books)[0].PublicationYear, uint32(2020))
	assert.Equal(t, len((*books)[0].Authors), 2)

}

func TestUpdateBook(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneBook()
	if err != nil {
		log.Fatal(err)
	}

	book, _ := models.ShowBook(app.DB, 1)

	assert.Equal(t, book.Name, "Book 1")
	assert.Equal(t, book.Edition, uint32(1))
	assert.Equal(t, book.PublicationYear, uint32(2020))
	assert.Equal(t, len(book.Authors), 2)

	newAuthors := []models.Author{
		{
			Name: "José Doe",
		},
	}

	for i := range newAuthors {
		app.DB.Model(&models.Author{}).Save(&newAuthors[i])
	}

	newAuthorsCreated, _ := models.GetAuthorsById(app.DB, []int{3})

	newBook := &models.Book{
		Name:            "Updated Book 1",
		Edition:         uint32(2),
		PublicationYear: uint32(2020),
		Authors:         *newAuthorsCreated,
	}

	updatedBook, err := newBook.UpdateBook(app.DB, int(book.ID))
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, updatedBook.Name, "Updated Book 1")
	assert.Equal(t, updatedBook.Edition, uint32(2))
	assert.Equal(t, updatedBook.PublicationYear, uint32(2020))
	assert.Equal(t, len(updatedBook.Authors), 1)
	assert.Equal(t, updatedBook.Authors[0].Name, "José Doe")
}
