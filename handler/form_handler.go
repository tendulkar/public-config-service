// handler/form_handler.go
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

type FormHandler struct {
	service *service.FormService
	logger  *slog.Logger
}

func NewFormHandler(service *service.FormService, logger *slog.Logger) *FormHandler {
	return &FormHandler{
		service: service,
		logger:  logger,
	}
}

func (h *FormHandler) GetAllForms(w http.ResponseWriter, r *http.Request) {
	forms, err := h.service.GetAllForms()
	if err != nil {
		h.logger.Error("error getting all forms", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forms)
}

func (h *FormHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	f, err := h.service.GetForm(id)
	if err != nil {
		h.logger.Error("error getting form", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if f == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func (h *FormHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	var f model.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateForm(&f); err != nil {
		h.logger.Error("error creating form", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

func (h *FormHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var f model.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		h.logger.Error("error decoding request body", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	f.ID = uint64(id)
	if err := h.service.UpdateForm(&f); err != nil {
		h.logger.Error("error updating form", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FormHandler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.logger.Error("error converting id", slog.Any("error", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteForm(id); err != nil {
		h.logger.Error("error deleting form", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
