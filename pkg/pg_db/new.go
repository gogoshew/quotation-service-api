package pg_db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	rateUpdateTable = ""
	latestRateTable = ""
)

type DatabasePg struct {
	DSN string
	db  *gorm.DB
}

func NewDatabasePg(dsn string) (*DatabasePg, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&RateUpdate{}, &LatestRate{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
		return nil, err
	}
	return &DatabasePg{
		DSN: dsn,
		db:  db,
	}, nil
}

func (d *DatabasePg) CreateRowForUpdate(rateUpdate *RateUpdate) error {
	res := d.db.Table(rateUpdateTable).Create(rateUpdate)
	if err := res.Error; err != nil {
		log.Fatalf("Error inserting row into %s table", rateUpdateTable)
		return err
	}
	return nil
}

func (d *DatabasePg) GetRowsForUpdate() {
	var ratesUpdate []RateUpdate
	d.db.Table(rateUpdateTable).Find(&ratesUpdate)
	if len(ratesUpdate) == 0 {
		log.Println("There is no rows for update")
		return nil
	}
	return ratesUpdate
}

func (d *DatabasePg) Bla() {

}
