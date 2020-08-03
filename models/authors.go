package models

import (
	"github.com/jinzhu/gorm"
)

type Author struct {
	gorm.Model
	Name string `gorm:"size:100;not null" json:"name"`
}

func (a *Author) Save(db *gorm.DB) (*Author, error) {
	err := db.Debug().Save(&a).Error
	if err != nil {
		return &Author{}, err
	}
	return a, nil
}

func GetAllAuthors(db *gorm.DB) (*[]Author, error) {
	authors := []Author{}

	if err := db.Debug().Table("authors").Find(&authors).Error; err != nil {
		return &[]Author{}, err
	}

	return &authors, nil
}

func GetAuthorsById(db *gorm.DB, ids []int) (*[]Author, error) {
	authors := []Author{}

	if err := db.Debug().Table("authors").Where(ids).Find(&authors).Error; err != nil {
		return &[]Author{}, err
	}

	return &authors, nil
}
