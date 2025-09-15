package currency_conversion

import (
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var Instance Interface

type CurrencyRate struct {
	Rate     string
	Time     time.Time
	Currency types.Currency
}

type Interface interface {
	GetLatestRates(base types.Currency, quotes []types.Currency) ([]CurrencyRate, error)
	SetupMetrics(reg *prometheus.Registry)
}
