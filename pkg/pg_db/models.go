package pg_db

import (
	"time"

	"gorm.io/gorm"
)

type RateUpdate struct {
	ID         uint   `gorm:"primarykey"`
	UpdateID   string `gorm:"unique"`
	Currency   string
	Value      float64
	CreateAt   time.Time
	UpdateDate time.Time
	DeletedAt  gorm.DeletedAt
}

type LatestRate struct {
	gorm.Model
	Currency string `gorm:"unique"`
	Value    float64
	Base     string
	Date     time.Time
}
