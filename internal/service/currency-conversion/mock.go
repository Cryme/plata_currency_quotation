package currency_conversion

import (
	"fmt"
	"math/rand"
	"plata_currency_quotation/internal/domain/types"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) SetupMetrics(_ *prometheus.Registry) {

}

func (m *Mock) GetLatestRates(_ types.Currency, quotes []types.Currency) ([]CurrencyRate, error) {
	var rates []CurrencyRate

	for _, quote := range quotes {
		rates = append(rates, CurrencyRate{Rate: fmt.Sprintf("%.2f", 0.1+rand.Float64()*2000.), Time: time.Now(), Currency: quote})
	}

	return rates, nil
}
