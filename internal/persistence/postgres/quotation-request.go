package postgres

import (
	"errors"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *Db) QuotationRequestCreateOrGetByIdempotencyKey(request *qr.QuotationRequest) error {
	result := d.inner.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idempotency_key"}},
		DoNothing: true,
	}).Create(request)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		request.Id = uuid.Nil

		//WTF: какой гений решил, что это хорошая идея неявно id добавлять к where...
		return d.inner.First(request, "idempotency_key = ?", request.IdempotencyKey).Error
	}

	return nil
}

func (d *Db) QuotationRequestGetById(id uuid.UUID) (*qr.QuotationRequest, error) {
	var request qr.QuotationRequest

	if err := d.inner.First(&request, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &qr.QuotationRequest{}, err
	}

	return &request, nil
}

func (d *Db) QuotationRequestUpdateByBaseAndQuote(baseCurrency types.Currency, quoteCurrency types.Currency, price string, time time.Time) error {
	return d.inner.Model(&qr.QuotationRequest{}).
		Where("base_currency = ? AND quote_currency = ?", baseCurrency, quoteCurrency).
		Update("rate", price).
		Update("completed_at", time).
		Error
}

func (d *Db) QuotationRequestGetUniqUnhandled() ([][2]types.Currency, error) {
	result := make([][2]types.Currency, 0)

	rows, err := d.inner.Model(&qr.QuotationRequest{}).
		Select("DISTINCT base_currency, quote_currency").
		Where("completed_at is null").
		Rows()

	if err != nil {
		return nil, err
	}

	defer func() {
		//TODO: не уверен, мб стоит просто залогировать
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		var baseCurrency, quoteCurrency types.Currency

		if err := rows.Scan(&baseCurrency, &quoteCurrency); err != nil {
			return nil, err
		}

		result = append(result, [2]types.Currency{baseCurrency, quoteCurrency})
	}

	return result, nil
}
