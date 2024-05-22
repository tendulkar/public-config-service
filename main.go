// main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"

	"stellarsky.ai/platform/public-config-service/config"
	"stellarsky.ai/platform/public-config-service/db"
	"stellarsky.ai/platform/public-config-service/handler"
	"stellarsky.ai/platform/public-config-service/middleware"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
	"stellarsky.ai/platform/public-config-service/service"
)

func main() {
	// Initialize Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load configuration
	cfg := config.LoadConfig(logger)

	// Initialize Database
	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Error("could not initialize database", slog.Any("error", err))
		return
	}

	// Automigrate models
	database.AutoMigrate(&model.Type{}, &model.Validation{}, &model.Attribute{}, &model.Form{})

	// Initialize Repositories
	typeRepo := repository.NewTypeRepository(database, logger)
	validationRepo := repository.NewValidationRepository(database, logger)
	attributeRepo := repository.NewAttributeRepository(database, logger)
	formRepo := repository.NewFormRepository(database, logger)

	// Initialize Services
	typeService := service.NewTypeService(typeRepo, logger)
	validationService := service.NewValidationService(validationRepo, logger)
	attributeService := service.NewAttributeService(attributeRepo, logger)
	formService := service.NewFormService(formRepo, logger)

	// Initialize Handlers
	typeHandler := handler.NewTypeHandler(typeService, logger)
	validationHandler := handler.NewValidationHandler(validationService, logger)
	attributeHandler := handler.NewAttributeHandler(attributeService, logger)
	formHandler := handler.NewFormHandler(formService, logger)

	// Initialize Router
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.TracingMiddleware)

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/types", typeHandler.GetAllTypes).Methods("GET")
	api.HandleFunc("/types", typeHandler.CreateType).Methods("POST")
	api.HandleFunc("/types/{id}", typeHandler.GetType).Methods("GET")
	api.HandleFunc("/types/{id}", typeHandler.UpdateType).Methods("PUT")
	api.HandleFunc("/types/{id}", typeHandler.DeleteType).Methods("DELETE")

	api.HandleFunc("/validations", validationHandler.GetAllValidations).Methods("GET")
	api.HandleFunc("/validations", validationHandler.CreateValidation).Methods("POST")
	api.HandleFunc("/validations/{id}", validationHandler.GetValidation).Methods("GET")
	api.HandleFunc("/validations/{id}", validationHandler.UpdateValidation).Methods("PUT")
	api.HandleFunc("/validations/{id}", validationHandler.DeleteValidation).Methods("DELETE")

	api.HandleFunc("/attributes", attributeHandler.GetAllAttributes).Methods("GET")
	api.HandleFunc("/attributes", attributeHandler.CreateAttribute).Methods("POST")
	api.HandleFunc("/attributes/{id}", attributeHandler.GetAttribute).Methods("GET")
	api.HandleFunc("/attributes/{id}", attributeHandler.UpdateAttribute).Methods("PUT")
	api.HandleFunc("/attributes/{id}", attributeHandler.DeleteAttribute).Methods("DELETE")

	api.HandleFunc("/forms", formHandler.GetAllForms).Methods("GET")
	api.HandleFunc("/forms", formHandler.CreateForm).Methods("POST")
	api.HandleFunc("/forms/{id}", formHandler.GetForm).Methods("GET")
	api.HandleFunc("/forms/{id}", formHandler.UpdateForm).Methods("PUT")
	api.HandleFunc("/forms/{id}", formHandler.DeleteForm).Methods("DELETE")

	// Initialize server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + cfg.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen error", slog.Any("error", err))
		}
	}()
	logger.Info("Server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown Failed", slog.Any("error", err))
	}
	logger.Info("Server Exited Properly")
}
