package api

import (
	"log/slog"
	"net/http"
	"plata_currency_quotation/internal/api/quotation"
	"plata_currency_quotation/internal/lib/config"
	"plata_currency_quotation/internal/lib/env"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"

	_ "plata_currency_quotation/docs"
)

const (
	SwaggerEndpoint string = "/docs"
)

// @title Plata Test Quotation API
// @version 0.1
// @description Bla bla

func RegisterRoutes(router *chi.Mux, log *slog.Logger) {
	router.Route("/api", func(router chi.Router) {
		quotation.RegisterRoutes(router, log)
	})

	if config.V.Env != env.Prod {
		handler := httpSwagger.WrapHandler

		if config.V.Env != env.Local {
			handler = basicAuth(handler)
		}

		router.Get(SwaggerEndpoint+"/*", handler)
	}
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != config.V.SwaggerUser || pass != config.V.SwaggerPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}
		handler(w, r)
	}
}
