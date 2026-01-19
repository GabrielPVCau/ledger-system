package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gabrielcau/ledger-system/internal/domain"
	"github.com/gabrielcau/ledger-system/internal/service"
)

type TransferHandler struct {
	service *service.LedgerService
}

func NewTransferHandler(service *service.LedgerService) *TransferHandler {
	return &TransferHandler{service: service}
}

type TransferRequest struct {
	From   int   `json:"from"`
	To     int   `json:"to"`
	Amount int64 `json:"amount"`
}

func (h *TransferHandler) MakeTransfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Transfer(r.Context(), req.From, req.To, req.Amount); err != nil {
		switch {
		case errors.Is(err, domain.ErrInsufficientFunds):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrAccountNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domain.ErrInvalidAmount), errors.Is(err, domain.ErrSameAccount):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "transfer successful"})
}
