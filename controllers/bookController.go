package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Lgdev07/olist_challenge/models"
	"github.com/Lgdev07/olist_challenge/utils"
	"github.com/gorilla/mux"
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

	authorsList, _ := models.GetAuthorsById(s.DB, response.Authors)

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

func (s *Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success"}

	params := mux.Vars(r)

	intID, _ := strconv.Atoi(params["id"])

	book, _ := models.ShowBook(s.DB, intID)
	if book.ID == 0 {
		resp["status"] = "failed"
		resp["message"] = "Book not Found"
	}

	if err := models.DeleteBook(s.DB, intID); err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["message"] = "Deleted Book " + book.Name
	utils.JSON(w, http.StatusOK, resp)
	return
}

func (s *Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Book Updated"}

	params := mux.Vars(r)

	intID, _ := strconv.Atoi(params["id"])

	response := &Response{}

	err := json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authorsList, _ := models.GetAuthorsById(s.DB, response.Authors)

	book := &models.Book{
		Name:            response.Name,
		Edition:         response.Edition,
		PublicationYear: response.PublicationYear,
		Authors:         *authorsList,
	}

	updatedBook, err := book.UpdateBook(s.DB, intID)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["book"] = updatedBook
	utils.JSON(w, http.StatusOK, resp)
	return
}

func (s *Server) ListBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	booksModel := &[]models.Book{}

	name := strings.Join(query["name"], "")
	publicationYear := strings.Join(query["publication_year"], "")
	edition := strings.Join(query["edition"], "")
	author := strings.Join(query["author"], "")

	chain := s.DB.Debug().Preload("Authors")
	chain = chain.Joins("inner join book_authors on book_authors.book_id = books.id")
	chain = chain.Joins("inner join authors on authors.id = book_authors.author_id")
	chain = chain.Where("")

	if name != "" {
		chain = chain.Where("books.name = ?", name)
	}

	if publicationYear != "" {
		chain = chain.Where("publication_year = ?", publicationYear)
	}

	if edition != "" {
		chain = chain.Where("edition = ?", edition)
	}

	if author != "" {
		chain = chain.Where("authors.name = ?", author)
	}

	err := chain.Group("books.id").Find(booksModel).Error
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	utils.JSON(w, http.StatusOK, booksModel)
	return

}
