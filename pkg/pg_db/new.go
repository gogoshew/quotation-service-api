package pg_db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	rateBufferTable = "buffer_rates"
	latestTable     = "latest_rates"
	baseCurrency    = "USD"
	eur             = "EUR"
	mxn             = "MXN"
	gel             = "GEL"
)

type DatabasePg struct {
	dsn string
	db  *gorm.DB
}

func NewDatabasePg(dsn string) (*DatabasePg, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&BufferRate{}, &LatestRate{})
	if err != nil {
		log.Printf("Failed to migrate: %v", err)
		return nil, err
	}

	latestRates := []LatestRate{
		{Currency: eur, Base: baseCurrency},
		{Currency: mxn, Base: baseCurrency},
		{Currency: gel, Base: baseCurrency},
	}

	var count int64
	db.Model(&LatestRate{}).Count(&count)
	if count == 0 {
		db.Create(&latestRates)
		if db.Error != nil {
			log.Printf("Error creating currencies for latest table: %v\n", err)
			return nil, err
		}
	}

	bufferRates := []BufferRate{
		{Currency: eur, Base: baseCurrency},
		{Currency: mxn, Base: baseCurrency},
		{Currency: gel, Base: baseCurrency},
	}
	var c int64
	db.Model(&BufferRate{}).Count(&c)
	if c == 0 {
		db.Create(&bufferRates)
		if db.Error != nil {
			log.Printf("Error creating rows for buffer rates table: %v\n", err)
			return nil, err
		}
	}

	return &DatabasePg{
		dsn: dsn,
		db:  db,
	}, nil
}

func (d *DatabasePg) UpdateBuffer(cur string, rateBuffer BufferRate) error {
	if err := d.db.Model(&BufferRate{}).Where("currency = ?", cur).Updates(rateBuffer).Error; err != nil {
		log.Printf("Error updating %s row in buffer_rates table: %v\n", cur, err)
		return err
	}
	return nil
}

func (d *DatabasePg) GetRowByCurBuffer(currency string) (BufferRate, error) {
	var row BufferRate

	if err := d.db.Table(rateBufferTable).Where("currency = ?", currency).Find(&row).Error; err != nil {
		return BufferRate{}, err
	}
	return row, nil
}

func (d *DatabasePg) GetRowByIDBuffer(updateID string) (BufferRate, error) {
	var row BufferRate

	if err := d.db.Table(rateBufferTable).Where("update_id = ?", updateID).Where("update_flag = ?", true).Find(&row).Error; err != nil {
		return BufferRate{}, err
	}
	return row, nil
}

func (d *DatabasePg) GetRowsForUpdate() ([]BufferRate, error) {
	var rows []BufferRate
	if err := d.db.Table(rateBufferTable).Where("update_flag = ?", true).Find(&rows).Error; err != nil {
		log.Printf("Error getting rows for update from %s table", rateBufferTable)
		return nil, err
	}

	return rows, nil
}

func (d *DatabasePg) UpdateLatest(currency string, val float64) error {
	if err := d.db.Model(&LatestRate{}).Where("currency = ?", currency).Update("value", val).Error; err != nil {
		log.Printf("Error updating latest table: %v\n", err)
		return err
	}
	if err := d.db.Model(&BufferRate{}).Where("currency = ?", currency).Update("update_flag", false).Error; err != nil {
		log.Printf("Error scaling to false update_flag in rate_buffer table: %v\n", err)
		return err
	}
	return nil
}

func (d *DatabasePg) GetLatest(currency string) (LatestRate, error) {
	var row LatestRate

	if err := d.db.Table(latestTable).Where("currency = ?", currency).Find(&row).Error; err != nil {
		return LatestRate{}, err
	}
	return row, nil
}
