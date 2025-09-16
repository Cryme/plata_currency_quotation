package postgres

import (
	"fmt"
	"log"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/lib/config"
	"plata_currency_quotation/internal/lib/logger/sl"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	inner *gorm.DB
}

func (d *Db) OnStart() error {
	if err := d.inner.AutoMigrate(&qr.QuotationRequest{}); err != nil {
		return err
	}

	return nil
}

func New() *Db {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.Instance.DbHost,
		config.Instance.DbUser,
		config.Instance.DbPassword,
		config.Instance.DbName,
		config.Instance.DbPort,
		config.Instance.DbSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to init postgres db", sl.Err(err))
	}

	var s = Db{
		inner: db,
	}

	if err := s.OnStart(); err != nil {
		log.Fatal("failed postgres OnStart", sl.Err(err))
	}

	return &s
}
