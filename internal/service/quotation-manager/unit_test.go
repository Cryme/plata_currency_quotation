package quotation_manager

import (
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/persistence"
	cc "plata_currency_quotation/internal/service/currency-conversion"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Runtime(t *testing.T) {
	db := persistence.NewInMemory()

	createAndAssert := func(base types.Currency, quote types.Currency) *qr.QuotationRequest {
		request, err := qr.New(base, quote, uuid.New())
		assert.NoError(t, err)

		err = db.QuotationRequestCreateOrGetByIdempotencyKey(&request)
		assert.NoError(t, err)

		return &request
	}

	request1 := createAndAssert(types.USD, types.MXN)
	request2 := createAndAssert(types.USD, types.MXN)
	request3 := createAndAssert(types.MXN, types.EUR)
	request4 := createAndAssert(types.EUR, types.MXN)

	manager := New(time.Second, db, cc.NewMock())

	manager.Run()

	time.Sleep(time.Second * 2)

	var assertUpdated = func(t *testing.T, request *qr.QuotationRequest) {
		requestUpdated, err := db.QuotationRequestGetById(request1.Id)

		assert.NoError(t, err)
		assert.NotNil(t, requestUpdated)
		assert.NotNil(t, requestUpdated.CompletedAt)
	}

	assertUpdated(t, request1)
	assertUpdated(t, request2)
	assertUpdated(t, request3)
	assertUpdated(t, request4)
}

func Test_GetQuotation(t *testing.T) {
	manager := New(time.Second, persistence.NewInMemory(), cc.NewMock())
	now := time.Now()

	manager.UpdateQuotation(types.USD, types.EUR, "1.5", now)

	info, exists := manager.GetQuotation(types.USD, types.EUR)
	if !exists {
		t.Error("Expected quotation to exist")
	}

	if info.Price != "1.5" {
		t.Errorf("Expected price 1.5, got %s", info.Price)
	}

	if !info.UpdatedAt.Equal(now) {
		t.Errorf("Expected updatedAt %v, got %v", now, info.UpdatedAt)
	}
}

func Test_UpdateQuotation(t *testing.T) {
	manager := New(time.Second, persistence.NewInMemory(), cc.NewMock())
	now := time.Now()

	manager.UpdateQuotation(types.USD, types.EUR, "1.5", now)

	assert.Len(t, manager.quotations, 1)

	quotation := manager.quotations[asKey(types.USD, types.EUR)]

	assert.NotNil(t, quotation)
	assert.Equal(t, "1.5", quotation.Price)
	assert.Equal(t, now, quotation.UpdatedAt)
}

func Test_GroupCurrencyPairs(t *testing.T) {
	pairs := [][2]types.Currency{
		{types.USD, types.EUR},
		{types.USD, types.MXN},
		{types.EUR, types.USD},
	}

	grouped := groupCurrencyPairs(pairs)

	expected := map[types.Currency][]types.Currency{
		types.USD: {types.EUR, types.MXN},
		types.EUR: {types.USD},
	}

	if !reflect.DeepEqual(grouped, expected) {
		t.Errorf("Expected %v, got %v", expected, grouped)
	}
}
