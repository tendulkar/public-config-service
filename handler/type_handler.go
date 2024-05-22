// handler/type_handler.go
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

type TypeHandler struct {
	service *service.TypeService
	logger  *slog.Logger
}

func NewTypeHandler(service *service.TypeService, logger *slog.Logger) *TypeHandler {
	return &TypeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TypeHandler) GetAllTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetAllTypes()
	if err != nil {
		h.logger.Error("error getting all types", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (h *TypeHandler) GetType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	t, err := h.service.GetType(id)
	if err != nil {
		h.logger.Error("error getting type", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TypeHandler) CreateType(w http.ResponseWriter, r *http.Request) {
	var t model.Type
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateType(&t); err != nil {
		h.logger.Error("error creating type", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TypeHandler) UpdateType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var t model.Type
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	t.ID = uint64(id)
	if err := h.service.UpdateType(&t); err != nil {
		h.logger.Error("error updating type", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TypeHandler) DeleteType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteType(id); err != nil {
		h.logger.Error("error deleting type", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
