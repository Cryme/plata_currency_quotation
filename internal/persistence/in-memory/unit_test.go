package in_memory

import (
	"testing"
	"time"

	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newTestDb() *Db {
	return &Db{store: make([]*qr.QuotationRequest, 0)}
}

func Test_CreateNewRequest(t *testing.T) {
	db := newTestDb()

	req := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	err := db.QuotationRequestCreateOrGetByIdempotencyKey(req)
	assert.NoError(t, err)
	assert.Len(t, db.store, 1)

	stored, err := db.QuotationRequestGetById(req.Id)
	assert.NoError(t, err)
	assert.NotNil(t, stored)
	assert.NotSame(t, stored, req)

	assert.Equal(t, stored.Id, req.Id)
}

func Test_CreateIdempotent(t *testing.T) {
	db := newTestDb()

	var key = uuid.New()

	req1 := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: key,
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	req2 := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: key,
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	err := db.QuotationRequestCreateOrGetByIdempotencyKey(req1)
	assert.NoError(t, err)

	err = db.QuotationRequestCreateOrGetByIdempotencyKey(req2)
	assert.NoError(t, err)

	assert.Len(t, db.store, 1)
	assert.Equal(t, req1.Id, req2.Id)
}

func Test_FoundAndNotFound(t *testing.T) {
	db := newTestDb()

	id := uuid.New()

	req := &qr.QuotationRequest{
		Id:             id,
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	err := db.QuotationRequestCreateOrGetByIdempotencyKey(req)
	assert.NoError(t, err)

	found, err := db.QuotationRequestGetById(id)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, id, found.Id)

	other, err := db.QuotationRequestGetById(uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, other)
}

func Test_UpdateByBaseAndQuote(t *testing.T) {
	db := newTestDb()
	now := time.Now()

	req := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	err := db.QuotationRequestCreateOrGetByIdempotencyKey(req)
	assert.NoError(t, err)

	err = db.QuotationRequestUpdateByBaseAndQuote(types.USD, types.EUR, "1.25", now)
	assert.NoError(t, err)

	req, err = db.QuotationRequestGetById(req.Id)
	assert.NoError(t, err)

	assert.NotNil(t, req)
	assert.NotNil(t, req.Rate)

	assert.Equal(t, "1.25", *req.Rate)
	assert.NotNil(t, req.CompletedAt)
	assert.WithinDuration(t, now, *req.CompletedAt, time.Second)
}

func Test_GetUniqUnhandled(t *testing.T) {
	db := newTestDb()

	req1 := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	req2 := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.EUR,
	}

	now := time.Now()
	rate := "2.4"
	req3 := &qr.QuotationRequest{
		Id:             uuid.New(),
		IdempotencyKey: uuid.New(),
		BaseCurrency:   types.USD,
		QuoteCurrency:  types.MXN,
		Rate:           &rate,
		CompletedAt:    &now,
	}

	err := db.QuotationRequestCreateOrGetByIdempotencyKey(req1)
	assert.NoError(t, err)

	err = db.QuotationRequestCreateOrGetByIdempotencyKey(req2)
	assert.NoError(t, err)

	err = db.QuotationRequestCreateOrGetByIdempotencyKey(req3)
	assert.NoError(t, err)

	keys, err := db.QuotationRequestGetUniqUnhandled()
	assert.NoError(t, err)

	assert.Equal(t, [][2]types.Currency{
		{types.USD, types.EUR},
	}, keys)
}
