package usecase

import (
	"context"
	"log/slog"
	"os"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/persistence"
	"plata_currency_quotation/internal/persistence/inmemory"
	cc "plata_currency_quotation/internal/service/currency-conversion"
	qm "plata_currency_quotation/internal/service/quotation-manager"
	"plata_currency_quotation/internal/usecase/command"
	qry "plata_currency_quotation/internal/usecase/query"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_RequestQuotationUpdate(t *testing.T) {
	persistence.Instance = inmemory.New()
	cc.Instance = cc.NewMock()

	qm.Instance = qm.New(
		time.Duration(10)*time.Millisecond,
		persistence.Instance,
		cc.Instance,
	)

	key := uuid.New()

	baseCurrency := types.USD
	quoteCurrency := types.MXN

	command := cmd.UpdateQuotation{BaseCurrency: baseCurrency, QuoteCurrency: quoteCurrency, IdempotencyKey: key}

	result, err := command.Execute(
		context.Background(),
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	assert.NoError(t, err)
	assert.NotEqual(t, result.Id, uuid.Nil)

	quotationRequest, err := persistence.Instance.QuotationRequestGetById(result.Id)

	assert.NoError(t, err)
	assert.NotNil(t, quotationRequest)
	assert.Equal(t, result.Id, quotationRequest.Id)
	assert.Equal(t, baseCurrency, quotationRequest.BaseCurrency)
	assert.Equal(t, quoteCurrency, quotationRequest.QuoteCurrency)
	assert.Nil(t, quotationRequest.CompletedAt)
	assert.Nil(t, quotationRequest.Rate)

	qm.Instance.Run()
	time.Sleep(time.Duration(100) * time.Millisecond)

	quotationRequest, err = persistence.Instance.QuotationRequestGetById(result.Id)

	assert.NoError(t, err)
	assert.NotNil(t, quotationRequest)
	assert.Equal(t, result.Id, quotationRequest.Id)
	assert.Equal(t, baseCurrency, quotationRequest.BaseCurrency)
	assert.Equal(t, quoteCurrency, quotationRequest.QuoteCurrency)
	assert.NotNil(t, quotationRequest.CompletedAt)
	assert.NotNil(t, quotationRequest.Rate)
}

func Test_RequestQuotationUpdateIdempotency(t *testing.T) {
	persistence.Instance = inmemory.New()
	cc.Instance = cc.NewMock()

	qm.Instance = qm.New(
		time.Duration(10)*time.Millisecond,
		persistence.Instance,
		cc.Instance,
	)

	key := uuid.New()

	baseCurrency := types.USD
	quoteCurrency := types.MXN

	command := cmd.UpdateQuotation{BaseCurrency: baseCurrency, QuoteCurrency: quoteCurrency, IdempotencyKey: key}

	result, err := command.Execute(
		context.Background(),
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	assert.NoError(t, err)
	assert.NotEqual(t, result.Id, uuid.Nil)

	firstId := result.Id

	command = cmd.UpdateQuotation{BaseCurrency: baseCurrency, QuoteCurrency: quoteCurrency, IdempotencyKey: key}

	result, err = command.Execute(
		context.Background(),
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	assert.NoError(t, err)

	assert.Equal(t, firstId, result.Id)
}

func Test_GetQuotationByRequestId(t *testing.T) {
	persistence.Instance = inmemory.New()
	cc.Instance = cc.NewMock()

	qm.Instance = qm.New(
		time.Duration(20)*time.Millisecond,
		persistence.Instance,
		cc.Instance,
	)

	command := cmd.UpdateQuotation{BaseCurrency: types.USD, QuoteCurrency: types.MXN, IdempotencyKey: uuid.New()}

	var id uuid.UUID

	{
		result, err := command.Execute(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.NoError(t, err)
		assert.NotEqual(t, result.Id, uuid.Nil)

		id = result.Id
	}

	{
		query := qry.GetQuotationByRequestId{Id: id}

		_, err := query.Run(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.Equal(t, err, qry.ErrRequestNotReady)
	}

	qm.Instance.Run()
	time.Sleep(time.Duration(100) * time.Millisecond)

	{
		query := qry.GetQuotationByRequestId{Id: id}

		result, err := query.Run(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEqual(t, result.UpdatedAt, 0)
		assert.NotEqual(t, result.Rate, "")
	}
}

func Test_GetQuotationByCurrency(t *testing.T) {
	persistence.Instance = inmemory.New()
	cc.Instance = cc.NewMock()

	qm.Instance = qm.New(
		time.Duration(20)*time.Millisecond,
		persistence.Instance,
		cc.Instance,
	)

	baseCurrency := types.USD
	quoteCurrency := types.MXN

	{
		query := qry.GetQuotation{Base: baseCurrency, Quote: quoteCurrency}

		_, err := query.Run(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.Equal(t, err, qry.ErrNoQuotationData)
	}

	command := cmd.UpdateQuotation{BaseCurrency: baseCurrency, QuoteCurrency: quoteCurrency, IdempotencyKey: uuid.New()}

	{
		result, err := command.Execute(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.NoError(t, err)
		assert.NotEqual(t, result.Id, uuid.Nil)
	}

	qm.Instance.Run()
	time.Sleep(time.Duration(100) * time.Millisecond)

	{
		query := qry.GetQuotation{Base: baseCurrency, Quote: quoteCurrency}

		result, err := query.Run(
			context.Background(),
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		)

		assert.NoError(t, err)
		assert.NotEqual(t, result.UpdatedAt, 0)
		assert.NotEqual(t, result.Rate, "")
	}
}
