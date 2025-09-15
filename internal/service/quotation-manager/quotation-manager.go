package quotation_manager

import (
	"log/slog"
	"os"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/lib/logger/sl"
	"plata_currency_quotation/internal/persistence"
	cc "plata_currency_quotation/internal/service/currency-conversion"
	"sync"
	"time"
)

var Instance *QuotationManager

type QuotationManager struct {
	runInterval     time.Duration
	mutex           sync.RWMutex
	quotations      map[string]types.QuotationInfo
	db              persistence.Interface
	currencyConvert cc.Interface
	logger          *slog.Logger
}

func New(runInterval time.Duration, db persistence.Interface, currencyConvert cc.Interface) *QuotationManager {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With(
		"component", "service/quotation-manager",
	)

	manager := QuotationManager{
		runInterval:     runInterval,
		quotations:      make(map[string]types.QuotationInfo),
		mutex:           sync.RWMutex{},
		db:              db,
		currencyConvert: currencyConvert,
		logger:          logger,
	}

	return &manager
}

func (q *QuotationManager) GetQuotation(base types.Currency, quote types.Currency) (types.QuotationInfo, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	info, exists := q.quotations[string(base+"/"+quote)]

	return info, exists
}

func (q *QuotationManager) UpdateQuotation(base types.Currency, quote types.Currency, price string, updatedAt time.Time) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.quotations[asKey(base, quote)] = types.QuotationInfo{
		Price:     price,
		UpdatedAt: updatedAt,
	}
}

func (q *QuotationManager) Run() {
	go func() {
		for {
			var startedAt = time.Now()

			q.runRequestsHandler()

			executionTime := time.Since(startedAt)
			sleepDuration := q.runInterval - executionTime

			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
		}
	}()
}

func asKey(base types.Currency, quote types.Currency) string {
	return string(base + "/" + quote)
}

func groupCurrencyPairs(pairs [][2]types.Currency) map[types.Currency][]types.Currency {
	grouped := make(map[types.Currency][]types.Currency)

	for _, pair := range pairs {
		base, quote := pair[0], pair[1]
		grouped[base] = append(grouped[base], quote)
	}

	return grouped
}

func (q *QuotationManager) runRequestsHandler() {
	currencyPairs, err := q.db.QuotationRequestGetUniqUnhandled()

	if err != nil {
		q.logger.Error("failed to get currency pairs", sl.Err(err))

		return
	}

	if len(currencyPairs) == 0 {
		return
	}

	groupedPairs := groupCurrencyPairs(currencyPairs)
	var wg sync.WaitGroup

	for base, quotes := range groupedPairs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			rates, err := q.currencyConvert.GetLatestRates(base, quotes)

			if err != nil {
				q.logger.Error("failed to get latest rates", sl.Err(err))
			}

			for _, rate := range rates {
				err = q.db.QuotationRequestUpdateByBaseAndQuote(base, rate.Currency, rate.Rate, rate.Time)

				if err != nil {
					q.logger.Error("failed to update quotation requests", sl.Err(err))

					continue
				}

				q.UpdateQuotation(base, rate.Currency, rate.Rate, rate.Time)
			}
		}()
	}

	wg.Wait()
}
