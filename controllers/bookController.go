package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Lgdev07/olist_challenge/models"
	"github.com/Lgdev07/olist_challenge/utils"
)

type Response struct {
	Name            string `json:"name"`
	Edition         uint32 `json:"edition"`
	PublicationYear uint32 `json:"publication_year"`
	Authors         []int  `json:"authors"`
}

func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Book created successfully"}

	response := &Response{}

	err := json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authors := &models.Author{}

	authorsList, _ := authors.GetAuthorsById(s.DB, response.Authors)

	book := &models.Book{
		Name:            response.Name,
		Edition:         response.Edition,
		PublicationYear: response.PublicationYear,
		Authors:         *authorsList,
	}

	bookCreated, err := book.Save(s.DB)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["book"] = bookCreated
	utils.JSON(w, http.StatusCreated, resp)
	return

}

func (s *Server) ListBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	name := strings.Join(query["name"], "")
	publication_year := strings.Join(query["publication_year"], "")
	edition := strings.Join(query["edition"], "")
	author := strings.Join(query["author"], "")

	var queryString string

	if name != "" {
		queryString = `"name = ?", ` + name
	}

	// AND
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	//// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

}
