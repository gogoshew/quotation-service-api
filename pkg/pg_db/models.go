package pg_db

import (
	"time"
)

type BufferRate struct {
	Currency   string `gorm:"primarykey"`
	UpdateID   string
	Value      float64
	Base       string
	UpdateFlag bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type LatestRate struct {
	Currency  string `gorm:"primarykey"`
	Value     float64
	Base      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
