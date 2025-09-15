package quotation_request

import (
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type QuotationRequest struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	IdempotencyKey uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"`
	CreatedAt      time.Time      `gorm:"type:timestamp;not null"`
	BaseCurrency   types.Currency `gorm:"type:varchar(3);not null"`
	QuoteCurrency  types.Currency `gorm:"type:varchar(3);not null"`
	CompletedAt    *time.Time     `gorm:"type:timestamp"`
	Rate           *string        `gorm:"type:text"`
}

func New(baseCurrency types.Currency, quoteCurrency types.Currency, idempotencyKey uuid.UUID) (QuotationRequest, error) {
	if baseCurrency == quoteCurrency {
		return QuotationRequest{}, ErrSameCurrency
	}

	return QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now(),
		BaseCurrency:   baseCurrency,
		QuoteCurrency:  quoteCurrency,
		CompletedAt:    nil,
		Rate:           nil,
	}, nil
}
