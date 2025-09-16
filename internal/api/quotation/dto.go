package quotation

import (
	"plata_currency_quotation/internal/domain/types"

	"github.com/google/uuid"
)

type RequestQuotationUpdateBody struct {
	BaseCurrency   types.Currency `json:"baseCurrency" validate:"required,enum"`
	QuoteCurrency  types.Currency `json:"quoteCurrency" validate:"required,enum"`
	IdempotencyKey uuid.UUID      `json:"idempotencyKey" format:"uuid" validate:"required,uuid"`
}

type RequestQuotationUpdateResponse struct {
	RequestId uuid.UUID `json:"requestId" swaggertype:"string" format:"uuid" binding:"required"`
}

type GetCurrencyListResponse struct {
	Currencies []types.Currency `json:"currencies" swaggertype:"array,string" format:"uuid"`
}

type RequestStatus string

const (
	Ready    RequestStatus = "Ready"
	NotReady RequestStatus = "NotReady"
)

// @Description fields `price` and `timestamp` are only presented when status is `Ready`
type GetQuotationByRequestIdResponse struct {
	Status RequestStatus `json:"status" binding:"required"`
	Rate   string        `json:"rate" example:"123.45" swaggertype:"string" format:"decimal"`
	// Unix timestamp in milliseconds
	UpdatedAt int64 `json:"updatedAt" example:"1694613600" swaggertype:"integer" format:"int64"`
}

type GetQuotationByRequestIdResponseNotReady struct {
	Status RequestStatus `json:"status"`
}

type GetQuotationResponse struct {
	Rate string `json:"rate" example:"123.45" swaggertype:"string" format:"decimal"`
	// Unix timestamp in milliseconds
	UpdatedAt int64 `json:"updatedAt" example:"1694613600" swaggertype:"integer" format:"int64"`
}
