package persistence

import (
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

// WTF: Утиные интерфейсы полная хрень!

type QuotationRequestPersistentOperations interface {
	QuotationRequestCreateOrGetByIdempotencyKey(*qr.QuotationRequest) error
	QuotationRequestGetById(id uuid.UUID) (*qr.QuotationRequest, error)
	QuotationRequestUpdateByBaseAndQuote(baseCurrency types.Currency, quoteCurrency types.Currency, rate string, time time.Time) error
	QuotationRequestGetUniqUnhandled() ([][2]types.Currency, error)
}
