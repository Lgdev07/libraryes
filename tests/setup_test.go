package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Lgdev07/olist_challenge/controllers"
	"github.com/Lgdev07/olist_challenge/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var app = controllers.Server{}
var bookInstance = models.Book{}
var authorInstance = models.Author{}

func TestMain(m *testing.M) {
	//Since we add our .env in .gitignore, Circle CI cannot see it, so see the else statement
	if _, err := os.Stat("./../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		Database()
	} else {
		CIBuild()
	}
	os.Exit(m.Run())
}

//When using CircleCI
func CIBuild() {
	var err error
	DBURL := fmt.Sprintf(`host=localhost port=5432 user=lgdev07 
	dbname=olist_challenge_test sslmode=disable password=docker`)
	app.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to postgres database\n")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the postgres database\n")
	}
}

func Database() {

	var err error
	DBURL := fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=disable 
	password=%s`, os.Getenv("TEST_DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("TEST_DB_NAME"),
		os.Getenv("DB_PASSWORD"))

	app.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to postgres database\n")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the postgres database\n")
	}

}

func refreshBookAndAuthorsTable() error {
	err := app.DB.DropTableIfExists(&models.Book{}, &models.Author{}).Error
	if err != nil {
		return err
	}
	err = app.DB.AutoMigrate(&models.Book{}, &models.Author{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables books and authors")
	return nil
}

func seedOneBook() (models.Book, error) {
	authors := []models.Author{
		{
			Name: "João Doe",
		},
		{
			Name: "Maria Doe",
		},
	}

	for i := range authors {
		app.DB.Model(&models.Author{}).Create(&authors[i])
	}

	var params map[string]string

	authorsCreated, _ := models.GetAllAuthors(app.DB, params)

	book := models.Book{
		Name:            "Book 1",
		Edition:         1,
		PublicationYear: 2020,
		Authors:         *authorsCreated,
	}

	err := app.DB.Model(&models.Book{}).Create(&book).Error
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func seedOneAuthor() error {
	author := models.Author{Name: "João Doe"}

	err := app.DB.Model(&models.Author{}).Create(&author).Error
	if err != nil {
		return err
	}
	return nil
}
