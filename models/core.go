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
	DeletedAt *time.Time
}

func init() {
	var err error

	db, err = gorm.Open("postgres", "host=127.0.0.1 user=byung dbname=BlogByung sslmode=disable password=1qaz@WSX")
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&User{})
}
