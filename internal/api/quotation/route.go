package quotation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/lib/http-server/response"
	"plata_currency_quotation/internal/lib/validator"
	"plata_currency_quotation/internal/usecase/command"
	qry "plata_currency_quotation/internal/usecase/query"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func RegisterRoutes(router chi.Router, log *slog.Logger) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/quotation", requestQuotationUpdate(log))
		router.Get("/quotation/{id}", getQuotationByRequestId(log))
		router.Get("/quotation/current", getQuotation(log))
	})
}

// @Summary Request quotation update
// @Description Creates a quotation update request
// @Tags Quotation
// @Accept json
// @Produce json
// @Param request body RequestQuotationUpdateBody true "Quotation request"
// @Success 200 {object} RequestQuotationUpdateResponse
// @Router /api/v1/quotation [post]
func requestQuotationUpdate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RequestQuotationUpdateBody

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := validator.Struct(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		command := cmd.UpdateQuotation{
			BaseCurrency:   request.BaseCurrency,
			QuoteCurrency:  request.QuoteCurrency,
			IdempotencyKey: request.IdempotencyKey,
		}

		result, err := command.Execute(r.Context(), log)

		if err != nil {
			switch {
			case errors.Is(err, qr.ErrSameCurrency):
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}

			return
		}

		response.Ok(w, log, RequestQuotationUpdateResponse{RequestId: result.Id})
	}
}

// @Summary Get quotation by ID
// @Description Retrieves a quotation by request ID. If request is not proceeded yet, returns status `NotReady`. If request is completed, returns status `Ready` and fields `rate` and `updatedAt`.
// @Tags Quotation
// @Produce json
// @Param id path string true "Quotation ID"
// @Success 200 {object} GetQuotationByRequestIdResponse
// @Router /api/v1/quotation/{id} [get]
func getQuotationByRequestId(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)

			return
		}

		query := qry.GetQuotationByRequestId{
			Id: id,
		}

		result, err := query.Run(r.Context(), log)

		if err != nil {
			switch {
			case errors.Is(err, qry.ErrNoRequestWithSuchId):
				http.Error(w, "no request with such id", http.StatusNotFound)
			case errors.Is(err, qry.ErrRequestNotReady):
				response.Ok(w, log, GetQuotationByRequestIdResponseNotReady{Status: NotReady})
			default:
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
			return
		}

		response.Ok(w, log, GetQuotationByRequestIdResponse{Rate: result.Rate, Status: Ready, UpdatedAt: result.UpdatedAt})
	}
}

// @Summary Get current quotation by currencies
// @Description Retrieves current quotation by base and quote currencies
// @Tags Quotation
// @Produce json
// @Param base query string true "Base Currency"
// @Param quote query string true "Quote Currency"
// @Success 200 {object} GetQuotationResponse
// @Router /api/v1/quotation/current [get]
func getQuotation(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := types.Currency(r.URL.Query().Get("base"))
		quote := types.Currency(r.URL.Query().Get("quote"))

		if !base.IsValid() {
			http.Error(w, "invalid base currency", http.StatusBadRequest)

			return
		}

		if !quote.IsValid() {
			http.Error(w, "invalid quote currency", http.StatusBadRequest)

			return
		}

		var query = qry.GetQuotation{
			Base:  base,
			Quote: quote,
		}

		var quotation, err = query.Run(r.Context(), log)

		if err != nil {
			switch {
			case errors.Is(err, qry.ErrQuotationNotFound):
				http.Error(w, "quotation not found", http.StatusNotFound)

				return
			default:
				http.Error(w, "something went wrong", http.StatusInternalServerError)

				return
			}
		}

		response.Ok(w, log, GetQuotationResponse{Rate: quotation.Price, UpdatedAt: quotation.UpdatedAt.UnixMilli()})
	}
}
