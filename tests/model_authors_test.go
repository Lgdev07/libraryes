package tests

import (
	"log"
	"testing"

	"github.com/Lgdev07/libraryes/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthor(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	if err = seedOneAuthor(); err != nil {
		log.Fatal(err)
	}

	var params map[string]string

	authors, err := models.GetAllAuthors(app.DB, params)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*authors), 1)
	assert.Equal(t, (*authors)[0].Name, "João Doe")
}

func TestGetByIdAuthor(t *testing.T) {
	err := refreshBookAndAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	if err = seedOneAuthor(); err != nil {
		log.Fatal(err)
	}

	authors, err := models.GetAuthorsById(app.DB, []int{2})
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*authors), 0)

	authors, err = models.GetAuthorsById(app.DB, []int{1})
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(*authors), 1)
	assert.Equal(t, (*authors)[0].Name, "João Doe")
}
