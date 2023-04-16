package db

import (
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

func Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("data"))
}
