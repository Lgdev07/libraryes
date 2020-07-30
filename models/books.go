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
