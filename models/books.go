package models

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Name            string   `gorm:"size:100;not null" json:"name"`
	Edition         uint32   `gorm:"not null" json:"edition"`
	PublicationYear uint32   `gorm:"not null" json:"publication_year"`
	Authors         []Author `gorm:"many2many:book_authors;" json:"authors"`
}

func (b *Book) Save(db *gorm.DB) (*Book, error) {
	err := db.Debug().Save(&b).Error
	if err != nil {
		return &Book{}, err
	}
	return b, nil

}

func (b *Book) GetAllBooks(db *gorm.DB) (*[]Book, error) {
	books := []Book{}

	err := db.Debug().Table("books").Find(&books).Error
	if err != nil {
		return &[]Book{}, err
	}
	return &books, nil
}

func ShowBook(db *gorm.DB, id int) (*Book, error) {
	book := &Book{}

	err := db.Debug().Table("books").Where("id = ?", id).First(book).Error
	if err != nil {
		return &Book{}, err
	}
	return book, nil
}

func DeleteBook(db *gorm.DB, id int) error {
	book := &Book{}

	err := db.Debug().Table("books").Where("id = ?", id).Delete(book).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *Book) UpdateBook(db *gorm.DB, id int) (*Book, error) {

	book := &Book{}

	db.Debug().Table("books").Where("id = ?", id).First(&book)

	book.Name = b.Name
	book.Edition = b.Edition
	book.PublicationYear = b.PublicationYear

	err := db.Save(&book).Association("Authors").Replace(b.Authors).Error
	if err != nil {
		return &Book{}, err
	}
	return book, nil
}
