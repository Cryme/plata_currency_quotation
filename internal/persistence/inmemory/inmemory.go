package inmemory

import (
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"sync"
)

type Db struct {
	store []*qr.QuotationRequest
	mutex sync.Mutex
}

func (d *Db) OnStart() error {
	return nil
}

func New() *Db {
	return &Db{
		store: make([]*qr.QuotationRequest, 0),
	}
}
