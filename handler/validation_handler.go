// handler/validation_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"

	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/service"
)

type ValidationHandler struct {
	service *service.ValidationService
	logger  *slog.Logger
}

func NewValidationHandler(service *service.ValidationService, logger *slog.Logger) *ValidationHandler {
	return &ValidationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ValidationHandler) GetAllValidations(w http.ResponseWriter, r *http.Request) {
	validations, err := h.service.GetAllValidations()
	if err != nil {
		h.logger.Error("error getting all validations", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validations)
}

func (h *ValidationHandler) GetValidation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	v, err := h.service.GetValidation(id)
	if err != nil {
		h.logger.Error("error getting validation", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if v == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *ValidationHandler) CreateValidation(w http.ResponseWriter, r *http.Request) {
	var v model.Validation
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateValidation(&v); err != nil {
		h.logger.Error("error creating validation", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func (h *ValidationHandler) UpdateValidation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var v model.Validation
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	v.ID = uint64(id)
	if err := h.service.UpdateValidation(&v); err != nil {
		h.logger.Error("error updating validation", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ValidationHandler) DeleteValidation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteValidation(id); err != nil {
		h.logger.Error("error deleting validation", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
