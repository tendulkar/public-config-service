// handler/attribute_handler.go
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

type AttributeHandler struct {
	service *service.AttributeService
	logger  *slog.Logger
}

func NewAttributeHandler(service *service.AttributeService, logger *slog.Logger) *AttributeHandler {
	return &AttributeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AttributeHandler) GetAllAttributes(w http.ResponseWriter, r *http.Request) {
	attributes, err := h.service.GetAllAttributes()
	if err != nil {
		h.logger.Error("error getting all attributes", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attributes)
}

func (h *AttributeHandler) GetAttribute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	a, err := h.service.GetAttribute(id)
	if err != nil {
		h.logger.Error("error getting attribute", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if a == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func (h *AttributeHandler) CreateAttribute(w http.ResponseWriter, r *http.Request) {
	var a model.Attribute
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateAttribute(&a); err != nil {
		h.logger.Error("error creating attribute", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func (h *AttributeHandler) UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var a model.Attribute
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	a.ID = uint64(id)
	if err := h.service.UpdateAttribute(&a); err != nil {
		h.logger.Error("error updating attribute", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AttributeHandler) DeleteAttribute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteAttribute(id); err != nil {
		h.logger.Error("error deleting attribute", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
