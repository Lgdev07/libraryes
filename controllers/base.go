package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Lgdev07/olist_challenge/middlewares"
	"github.com/Lgdev07/olist_challenge/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable 
		password=%s`, DbHost, DbPort, DbUser, DbName, DbPassword)

	s.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Println("Cannot conect to database %s", DbName)
		log.Fatal(err)
	} else {
		fmt.Println("We are connected to database %s", DbName)
	}

	s.DB.Debug().AutoMigrate(
		&models.Author{},
		&models.Book{},
	)

	s.Router = mux.NewRouter().StrictSlash(true)
	s.InitializeRoutes()
}

func (s *Server) InitializeRoutes() {
	s.Router.Use(middlewares.SetContentTypeMiddleware)

	s.Router.HandleFunc("/authors/import", s.ImportCsvAuthor).Methods("POST")
	s.Router.HandleFunc("/authors", s.CreateAuthor).Methods("POST")
	s.Router.HandleFunc("/authors", s.ListAuthors).Methods("GET")
	s.Router.HandleFunc("/books", s.createBook).Methods("POST")
	s.Router.HandleFunc("/books", s.ListBooks).Methods("GET")
	s.Router.HandleFunc("/books/{id:[0-9]+}", s.DeleteBook).Methods("DELETE")
	s.Router.HandleFunc("/books/{id:[0-9]+}", s.UpdateBook).Methods("PUT")
}

func (s *Server) RunServer() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("\nServer starting on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, s.Router))
}
