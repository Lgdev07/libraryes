package tests

import (
	"log"
	"testing"

	"github.com/Lgdev07/olist_challenge/models"
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

	authors, err := models.GetAllAuthors(app.DB)
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
