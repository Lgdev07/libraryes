package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Lgdev07/olist_challenge/models"
	"github.com/Lgdev07/olist_challenge/utils"
)

func (s *Server) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Author successfully created"}

	author := &models.Author{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &author)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authorCreated, err := author.Save(s.DB)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["author"] = authorCreated
	utils.JSON(w, http.StatusCreated, resp)
	return

}

func (s *Server) ListAuthors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := map[string]string{
		"name": strings.Join(query["name"], ""),
		"page": strings.Join(query["page"], ""),
	}

	authorsList, err := models.GetAllAuthors(s.DB, params)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}

	utils.JSON(w, http.StatusOK, authorsList)
	return
}

func (s *Server) ImportCsvAuthor(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Authors successfully created"}

	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])

	io.Copy(&buf, file)

	contents := buf.String()
	splitedContent := strings.Split(contents, ",\n")

	authors := []models.Author{}

	for line := range splitedContent {
		if line == 0 || splitedContent[line] == "" {
			continue
		}

		author := &models.Author{
			Name: splitedContent[line],
		}

		authorCreated, err := author.Save(s.DB)
		if err != nil {
			utils.ERROR(w, http.StatusBadRequest, err)
			return
		}

		authors = append(authors, *authorCreated)

	}

	buf.Reset()

	resp["authors"] = authors
	utils.JSON(w, http.StatusCreated, resp)
	return
}
