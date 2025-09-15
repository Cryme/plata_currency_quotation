package currency_conversion

import (
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CurrencyRate struct {
	Rate     string
	Time     time.Time
	Currency types.Currency
}

type Interface interface {
	GetLatestRates(base types.Currency, quotes []types.Currency) ([]CurrencyRate, error)
	SetupMetrics(reg *prometheus.Registry)
}
