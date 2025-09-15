package currency_conversion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/lib/logger/sl"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outgoing_http_requests_total",
			Help: "Total number of outgoing HTTP requests",
		},
		[]string{"method", "status", "service"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "outgoing_http_request_duration_seconds",
			Help:    "Duration of outgoing HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "service"},
	)
)

type FrankfurterApi struct {
	apiUrl string
	client *http.Client
}

func NewFrankfurterApi(apiUrl string) *FrankfurterApi {
	return &FrankfurterApi{
		apiUrl: apiUrl,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

type frankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func (m *FrankfurterApi) SetupMetrics(reg *prometheus.Registry) {
	reg.MustRegister(requestsTotal)
	reg.MustRegister(requestDuration)
}

func (m *FrankfurterApi) GetLatestRates(base types.Currency, quotes []types.Currency) ([]CurrencyRate, error) {
	symbols := ""
	for i, quote := range quotes {
		if i > 0 {
			symbols += ","
		}
		symbols += string(quote)
	}

	url := fmt.Sprintf("%s/v1/latest?base=%s&symbols=%s", m.apiUrl, base, symbols)

	start := time.Now()
	resp, err := m.client.Get(url)
	duration := time.Since(start).Seconds()

	requestDuration.WithLabelValues("GET", "frankfurter").Observe(duration)

	if err != nil {
		requestsTotal.WithLabelValues("GET", "error", "frankfurter").Inc()
		return nil, fmt.Errorf("failed to fetch rate: %w", err)
	}

	defer func() {
		err := resp.Body.Close()

		if err != nil {
			sl.Log.Error("failed to close response body", sl.Err(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		requestsTotal.WithLabelValues("GET", fmt.Sprint(resp.StatusCode), "frankfurter").Inc()

		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	requestsTotal.WithLabelValues("GET", fmt.Sprint(resp.StatusCode), "frankfurter").Inc()

	var response frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var rates []CurrencyRate

	for _, quote := range quotes {
		rate, exists := response.Rates[string(quote)]

		if !exists {
			//TODO: не уверен, мб стоит просто залогировать
			return nil, fmt.Errorf("rate not found for currency pair %s/%s", base, quote)
		}

		rates = append(rates, CurrencyRate{
			Rate: fmt.Sprintf("%g", rate),
			//TODO: Апи отдает только дату, без времени
			Time:     time.Now(),
			Currency: quote,
		})
	}

	return rates, nil
}
