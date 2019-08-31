package models

import (
	"time"

	_ "github.com/bmizerany/pq"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-"`
}

func init() {
	var err error

	db, err = gorm.Open("postgres", "host=127.0.0.1 user=byung dbname=blogbyung sslmode=disable password=1qaz@WSX")
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Image{})
}
