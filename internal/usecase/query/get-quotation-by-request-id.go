package qry

import (
	"context"
	"errors"
	"log/slog"
	gs "plata_currency_quotation/internal"
	"plata_currency_quotation/internal/lib/logger/sl"

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

func (q *GetQuotationByRequestId) Run(ctx context.Context, log *slog.Logger) (GetQuotationByRequestIdResponse, error) {
	quotationRequest, err := gs.Db.QuotationRequestGetById(q.Id)

	if err != nil {
		log.Error("failed to get quotation request", sl.Err(err), sl.TraceId(ctx))

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
