package pg_db

import (
	"time"
)

type RateBuffer struct {
	ID         uint `gorm:"primarykey"`
	UpdateID   string
	Currency   string `gorm:"unique"`
	Value      float64
	Base       string
	UpdateFlag bool
	CreateAt   time.Time
	UpdateAt   time.Time
}

type LatestRate struct {
	ID       uint   `gorm:"primarykey"`
	Currency string `gorm:"unique"`
	Value    float64
	Base     string
	CreateAt time.Time
	UpdateAt time.Time `gorm:"autoUpdateTime:false"`
}

/*
1. Приходит запрос на обновление котировки
2. Создается запрос в сервис с актуальной инфой по котировкам
3. Полученная с внешнего сервиса инфа должна записывается в таблицу RateUpdate,
	- Если по указанной валюте есть запись с флагом UpdateFlag = true, то мы берем UpdateID из Таблицы в БД и

*/
