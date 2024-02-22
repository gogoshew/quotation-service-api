package pg_db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	rateBufferTable = "rate_buffer"
	latestTable     = "latest"
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
		log.Fatalf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&RateBuffer{}, &LatestRate{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
		return nil, err
	}

	latestRates := []LatestRate{
		{Currency: eur, Base: baseCurrency},
		{Currency: mxn, Base: baseCurrency},
		{Currency: gel, Base: baseCurrency},
	}
	db.Create(&latestRates)
	if db.Error != nil {
		log.Fatalf("Error creating currencies for latest table: %v\n", err)
		return nil, err
	}

	return &DatabasePg{
		dsn: dsn,
		db:  db,
	}, nil
}

func (d *DatabasePg) CreateRowBuffer(rateBuffer *RateBuffer) error {
	res := d.db.Table(rateBufferTable).Create(rateBuffer)
	if err := res.Error; err != nil {
		log.Fatalf("Error inserting row into %s table", rateBufferTable)
		return err
	}
	return nil
}

func (d *DatabasePg) GetRowByCurBuffer(currency string) (RateBuffer, error) {
	var row RateBuffer

	if err := d.db.Table(rateBufferTable).Where("currency = ?", currency).Find(&row).Error; err != nil {
		return RateBuffer{}, err
	}
	return row, nil
}

func (d *DatabasePg) GetRowByIDBuffer(updateID string) (RateBuffer, error) {
	var row RateBuffer

	if err := d.db.Table(rateBufferTable).Where("update_id = ?", updateID).Find(&row).Error; err != nil {
		return RateBuffer{}, err
	}
	return row, nil
}

func (d *DatabasePg) GetRowsForUpdate() ([]RateBuffer, error) {
	var rows []RateBuffer
	if err := d.db.Table(rateBufferTable).Where("update_flag = ?", true).Find(&rows).Error; err != nil {
		log.Fatalf("Error getting rows for update from %s table", rateBufferTable)
		return nil, err
	}

	return rows, nil
}

func (d *DatabasePg) UpdateLatest(currency string, val float64) error {
	if err := d.db.Model(&LatestRate{}).Where("currency = ?", currency).Update("value", val).Error; err != nil {
		log.Fatalf("Error updating latest table: %v\n", err)
		return err
	}
	if err := d.db.Model(&RateBuffer{}).Where("currency = ?", currency).Update("update_flag", false).Error; err != nil {
		log.Fatalf("Error scaling to false update_flag in rate_buffer table: %v\n", err)
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
