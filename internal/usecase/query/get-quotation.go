package qry

import (
	"context"
	"errors"
	"log/slog"
	quotation_request "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	qm "plata_currency_quotation/internal/service/quotation-manager"
)

var ErrNoQuotationData = errors.New("quotation was not requested yet")

type GetQuotation struct {
	Base  types.Currency
	Quote types.Currency
}

type GetQuotationResponse struct {
	Rate      string
	UpdatedAt int64
}

func (q *GetQuotation) Run(_ context.Context, _ *slog.Logger) (types.QuotationInfo, error) {
	if q.Quote == q.Base {
		return types.QuotationInfo{}, quotation_request.ErrSameCurrency
	}

	quotation, found := qm.Instance.GetQuotation(q.Base, q.Quote)

	if !found {
		return types.QuotationInfo{}, ErrNoQuotationData
	}

	return quotation, nil
}
