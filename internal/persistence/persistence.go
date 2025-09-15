package persistence

import (
	"log"
	"plata_currency_quotation/internal/lib/logger/sl"
	in_memory "plata_currency_quotation/internal/persistence/in-memory"
	"plata_currency_quotation/internal/persistence/postgres"
)

type CommonPersistenceOperations interface {
	OnStart() error
}

type Interface interface {
	CommonPersistenceOperations
	QuotationRequestPersistentOperations
}

func New() Interface {
	db, err := postgres.New()

	if err != nil {
		log.Fatal("failed to init db", sl.Err(err))
	}

	if err = db.OnStart(); err != nil {
		log.Fatal("failed to start db", sl.Err(err))
	}

	return db
}

func NewInMemory() Interface {
	return in_memory.New()
}
