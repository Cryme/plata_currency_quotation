package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"plata_currency_quotation/internal/lib/logger/sl"
)

func Ok(w http.ResponseWriter, log *slog.Logger, body any) {

	if body == nil {
		w.WriteHeader(http.StatusOK)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Error("failed to send response", sl.Err(err))
	}

}

type ErrorResponse struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(ErrorResponse{Message: msg}); err != nil {
		log.Error("failed to send response", sl.Err(err))
	}
}
