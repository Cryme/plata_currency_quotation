package cmd

import (
	"context"
	"log/slog"
	gs "plata_currency_quotation/internal"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/lib/logger/sl"

	"github.com/google/uuid"
)

type UpdateQuotation struct {
	BaseCurrency   types.Currency
	QuoteCurrency  types.Currency
	IdempotencyKey uuid.UUID
}

type Result struct {
	Id uuid.UUID
}

func (u UpdateQuotation) Execute(ctx context.Context, log *slog.Logger) (Result, error) {
	quotationRequest, err := qr.New(u.BaseCurrency, u.QuoteCurrency, u.IdempotencyKey)

	if err != nil {
		log.Error("failed to create quotation request", sl.Err(err), sl.TraceId(ctx))

		return Result{}, err
	}

	err = gs.Db.QuotationRequestCreateOrGetByIdempotencyKey(&quotationRequest)

	if err != nil {
		log.Error("failed to create quotation request in db", sl.Err(err), sl.TraceId(ctx))

		return Result{}, err
	}

	return Result{
		Id: quotationRequest.Id,
	}, nil
}
