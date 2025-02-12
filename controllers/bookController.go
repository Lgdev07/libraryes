package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Lgdev07/libraryes/models"
	"github.com/Lgdev07/libraryes/utils"
	"github.com/gorilla/mux"
)

type Response struct {
	Name            string `json:"name"`
	Edition         uint32 `json:"edition"`
	PublicationYear uint32 `json:"publication_year"`
	Authors         []int  `json:"authors"`
}

func (s *Server) CreateBook(w http.ResponseWriter, r *http.Request) {
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

	params := map[string]string{
		"name":            strings.Join(query["name"], ""),
		"publicationYear": strings.Join(query["publication_year"], ""),
		"edition":         strings.Join(query["edition"], ""),
		"author":          strings.Join(query["author"], ""),
	}

	books, err := models.GetAllBooks(s.DB, params)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	utils.JSON(w, http.StatusOK, books)
	return

}
