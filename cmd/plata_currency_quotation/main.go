package main

import (
	"log/slog"
	"net/http"
	"os"
	"plata_currency_quotation/internal/api"
	"plata_currency_quotation/internal/lib/config"
	"plata_currency_quotation/internal/lib/env"
	"plata_currency_quotation/internal/lib/http-server/middleware/logger"
	metricsMiddleware "plata_currency_quotation/internal/lib/http-server/middleware/metrics"
	"plata_currency_quotation/internal/lib/http-server/middleware/trace-id"
	"plata_currency_quotation/internal/lib/logger/sl"
	"plata_currency_quotation/internal/lib/metrics"
	"plata_currency_quotation/internal/lib/validator"
	"plata_currency_quotation/internal/persistence"
	"plata_currency_quotation/internal/persistence/postgres"
	cc "plata_currency_quotation/internal/service/currency-conversion"
	qm "plata_currency_quotation/internal/service/quotation-manager"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.V = config.FromEnv()

	validator.RegisterValidators()

	setupLogger()

	sl.Log.Info("starting server", slog.String("env", string(config.V.Env)))
	sl.Log.Debug("debug messages are enabled")

	setupServices()

	metrics.Run(config.V.MetricsPort, cc.Instance)

	router := chi.NewRouter()

	router.Use(trace_id.New())
	router.Use(metricsMiddleware.New())
	router.Use(logger.New(sl.Log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(config.V.IncomingRequestTimeout))

	api.RegisterRoutes(router, sl.Log)

	address := config.V.ServerIp + ":" + strconv.Itoa(int(config.V.ServerPort))

	if err := http.ListenAndServe(address, router); err != nil {
		sl.Log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupServices() {
	cc.Instance = cc.NewFrankfurterApi(config.V.FrankfurterApiUrl)

	persistence.Instance = postgres.New()

	qm.Instance = qm.New(
		time.Duration(config.V.QuotationUpdateIntervalMilliseconds)*time.Millisecond,
		persistence.Instance,
		cc.Instance,
	)

	qm.Instance.Run()
}

func setupLogger() {
	switch config.V.Env {
	case env.Local:
		sl.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.Dev:
		sl.Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.Preprod, env.Prod:
		sl.Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
