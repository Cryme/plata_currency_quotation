package postgres

import (
	"fmt"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/lib/config"

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

func New() (*Db, error) {
	const op = "storage.postgres.New"

	var useSsl string

	if config.V.DbUseSsl {
		useSsl = "enable"
	} else {
		useSsl = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.V.DbHost, config.V.DbUser, config.V.DbPassword, config.V.DbName, config.V.DbPort, useSsl)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var s = Db{
		inner: db,
	}

	if err := s.OnStart(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &s, nil
}
