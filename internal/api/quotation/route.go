package quotation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	qr "plata_currency_quotation/internal/domain/enity/quotation-request"
	"plata_currency_quotation/internal/domain/types"
	"plata_currency_quotation/internal/lib/http-server/response"
	"plata_currency_quotation/internal/lib/logger/sl"
	"plata_currency_quotation/internal/lib/validator"
	"plata_currency_quotation/internal/usecase/command"
	qry "plata_currency_quotation/internal/usecase/query"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func RegisterRoutes(router chi.Router, log *slog.Logger) {
	router.Route("/v1", func(router chi.Router) {
		router.Post("/quotation/update-request", requestQuotationUpdate(log))
		router.Get("/quotation/update-request/{id}", getQuotationByRequestId(log))
		router.Get("/quotation/last-requested", getQuotation(log))
		router.Get("/currency/list", getCurrencyList(log))
	})
}

// @Summary Get list of supported currencies
// @Description Returns list of supported currency codes [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217)
// @Tags Currency
// @Produce json
// @Success 200 {object} GetCurrencyListResponse
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/currency/list [get]
func getCurrencyList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(sl.TraceId(r.Context()))

		response.Ok(w, log, GetCurrencyListResponse{Currencies: types.AllCurrencies()})
	}
}

// @Summary Request quotation update
// @Description Creates a quotation update request. Use [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217) currency code. List of supported currencies - `GET /api/v1/currency/list`. Returns request Id.
// @Tags Quotation
// @Accept json
// @Produce json
// @Param request body RequestQuotationUpdateBody true "Quotation request"
// @Success 200 {object} RequestQuotationUpdateResponse
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/quotation/update-request [post]
func requestQuotationUpdate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RequestQuotationUpdateBody

		log = log.With(sl.TraceId(r.Context()))

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, err.Error(), log)

			return
		}

		if err := validator.Struct(request); err != nil {
			response.Error(w, http.StatusBadRequest, err.Error(), log)

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
				response.Error(w, http.StatusBadRequest, "Currencies can't be same", log)
			default:
				response.Error(w, http.StatusInternalServerError, "Something went wrong", log)
			}

			return
		}

		response.Ok(w, log, RequestQuotationUpdateResponse{RequestId: result.Id})
	}
}

// @Summary Get quotation by request Id
// @Description Retrieves a quotation by request Id. If request is not proceeded yet, returns status `NotReady`. If request is completed, returns status `Ready` and fields `rate` and `updatedAt`.
// @Tags Quotation
// @Produce json
// @Param id path string true "Quotation ID"
// @Success 200 {object} GetQuotationByRequestIdResponse
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 404 {object} response.ErrorResponse "No request with such id"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/quotation/update-request/{id} [get]
func getQuotationByRequestId(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		log = log.With(sl.TraceId(r.Context()))

		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id format. Should be uuid", log)

			return
		}

		query := qry.GetQuotationByRequestId{
			Id: id,
		}

		result, err := query.Run(r.Context(), log)

		if err != nil {
			switch {
			case errors.Is(err, qry.ErrNoRequestWithSuchId):
				response.Error(w, http.StatusNotFound, "No request with such id", log)
			case errors.Is(err, qry.ErrRequestNotReady):
				response.Ok(w, log, GetQuotationByRequestIdResponseNotReady{Status: NotReady})
			default:
				response.Error(w, http.StatusInternalServerError, "Something went wrong", log)
			}

			return
		}

		response.Ok(w, log, GetQuotationByRequestIdResponse{Rate: result.Rate, Status: Ready, UpdatedAt: result.UpdatedAt})
	}
}

// @Summary Get last requested quotation by currencies
// @Description Retrieves last requested quotation by base and quote currencies. Use [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217) currency code. List of supported currencies - `GET /api/v1/currency/list`. Returns `404 Quotation not found` if quotation wasn't requested at least once, use `POST /api/v1/update-request` in this case
// @Tags Quotation
// @Produce json
// @Param base query string true "Base Currency"
// @Param quote query string true "Quote Currency"
// @Success 200 {object} GetQuotationResponse
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 404 {object} response.ErrorResponse "Quotation not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/quotation/last-requested [get]
func getQuotation(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := types.Currency(r.URL.Query().Get("base"))
		quote := types.Currency(r.URL.Query().Get("quote"))

		log = log.With(sl.TraceId(r.Context()))

		if !base.IsValid() {
			response.Error(w, http.StatusBadRequest, "Invalid base currency", log)

			return
		}

		if !quote.IsValid() {
			response.Error(w, http.StatusBadRequest, "Invalid quote currency", log)

			return
		}

		var query = qry.GetQuotation{
			Base:  base,
			Quote: quote,
		}

		var quotation, err = query.Run(r.Context(), log)

		if err != nil {
			switch {
			case errors.Is(err, qr.ErrSameCurrency):
				response.Error(w, http.StatusBadRequest, "Currencies can't be same", log)
			case errors.Is(err, qry.ErrNoQuotationData):
				response.Error(w, http.StatusNotFound, "Quotation was not requested yet", log)
			default:
				response.Error(w, http.StatusInternalServerError, "Something went wrong", log)
			}

			return
		}

		response.Ok(w, log, GetQuotationResponse{Rate: quotation.Price, UpdatedAt: quotation.UpdatedAt.UnixMilli()})
	}
}
