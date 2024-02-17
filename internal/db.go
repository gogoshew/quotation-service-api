package internal

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	dsn := "user=username password=password dbname=dbname host=localhost port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type RateModel struct {
	Rate  string `gorm:"primaryKey"`
	Value float64
	Base  string
	Date  time.Time
}
