package qry

import (
	"context"
	"errors"
	"log/slog"
	"plata_currency_quotation/internal/lib/logger/sl"
	"plata_currency_quotation/internal/persistence"

	"github.com/google/uuid"
)

var ErrRequestNotReady = errors.New("not ready, try again later")

var ErrNoRequestWithSuchId = errors.New("no request with such id")

type GetQuotationByRequestId struct {
	Id uuid.UUID
}

type GetQuotationByRequestIdResponse struct {
	Rate      string
	UpdatedAt int64
}

func (q *GetQuotationByRequestId) Run(_ context.Context, log *slog.Logger) (GetQuotationByRequestIdResponse, error) {
	quotationRequest, err := persistence.Instance.QuotationRequestGetById(q.Id)

	if err != nil {
		log.Error("failed to get quotation request", sl.Err(err))

		return GetQuotationByRequestIdResponse{}, err
	}

	if quotationRequest == nil {
		return GetQuotationByRequestIdResponse{}, ErrNoRequestWithSuchId
	}

	if quotationRequest.CompletedAt == nil {
		return GetQuotationByRequestIdResponse{}, ErrRequestNotReady
	}

	return GetQuotationByRequestIdResponse{
		Rate:      *quotationRequest.Rate,
		UpdatedAt: quotationRequest.CompletedAt.UnixMilli(),
	}, nil
}
