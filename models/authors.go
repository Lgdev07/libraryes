package models

import (
	"strconv"

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

func GetAllAuthors(db *gorm.DB, params map[string]string) (*[]Author, error) {
	authors := []Author{}
	limit := 10

	chain := db.Debug().Table("authors")

	if params["name"] != "" {
		chain = chain.Where("name = ?", params["name"])
	}

	page, _ := strconv.Atoi(params["page"])

	if page == 0 {
		page = 1
	}

	page = page * limit
	offset := page - limit

	if err := chain.Limit(limit).Offset(offset).Find(&authors).Error; err != nil {
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
