package inmemory

import (
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

func (d *Db) QuotationRequestCreateOrGetByIdempotencyKey(request *qr.QuotationRequest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, existing := range d.store {
		if existing.IdempotencyKey == request.IdempotencyKey {
			deepClone(existing, request)

			return nil
		}
	}

	var clone qr.QuotationRequest

	deepClone(request, &clone)

	d.store = append(d.store, &clone)

	return nil
}

func (d *Db) QuotationRequestGetById(id uuid.UUID) (*qr.QuotationRequest, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, req := range d.store {
		if req.Id == id {
			var clone qr.QuotationRequest

			deepClone(req, &clone)

			return &clone, nil
		}
	}

	return nil, nil
}

func (d *Db) QuotationRequestUpdateByBaseAndQuote(baseCurrency types.Currency, quoteCurrency types.Currency, price string, t time.Time) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, req := range d.store {
		if req.BaseCurrency == baseCurrency && req.QuoteCurrency == quoteCurrency {
			req.Rate = &price
			req.CompletedAt = &t
		}
	}

	return nil
}

func (d *Db) QuotationRequestGetUniqUnhandled() ([][2]types.Currency, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	result := make([][2]types.Currency, 0)

outer:
	for _, req := range d.store {
		if req.CompletedAt == nil {
			key := [2]types.Currency{req.BaseCurrency, req.QuoteCurrency}

			for _, existing := range result {
				if existing == key {
					continue outer
				}
			}

			result = append(result, key)
		}
	}

	return result, nil
}

func deepClone(src *qr.QuotationRequest, dst *qr.QuotationRequest) {
	if src == nil {
		return
	}

	*dst = *src

	if src.CompletedAt != nil {
		t := *src.CompletedAt
		dst.CompletedAt = &t

	}

	if src.Rate != nil {
		r := *src.Rate
		dst.Rate = &r
	}
}
