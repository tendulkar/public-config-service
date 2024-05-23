package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"stellarsky.ai/platform/public-config-service/handler"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
	"stellarsky.ai/platform/public-config-service/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func convertToGinFunc(f func(http.ResponseWriter, *http.Request)) gin.HandlerFunc {

// 	return func(c *gin.Context) {
// 		f(c.Writer, c.Request)
// 	}
// }

func setupRouter(db *gorm.DB, logger *slog.Logger) *mux.Router {

	typeRepo := repository.NewTypeRepository(db, logger)
	validationRepo := repository.NewValidationRepository(db, logger)
	attributeRepo := repository.NewAttributeRepository(db, logger)
	formRepo := repository.NewFormRepository(db, logger)

	typeService := service.NewTypeService(typeRepo, logger)
	validationService := service.NewValidationService(validationRepo, logger)
	attributeService := service.NewAttributeService(attributeRepo, logger)
	formService := service.NewFormService(formRepo, logger)

	typeHandler := handler.NewTypeHandler(typeService, logger)
	validationHandler := handler.NewValidationHandler(validationService, logger)
	attributeHandler := handler.NewAttributeHandler(attributeService, logger)
	formHandler := handler.NewFormHandler(formService, logger)

	// Routes
	// Type Routes
	// Validation Routes
	// Attribute Routes
	// Form Routes
	// r := setupGinRouter(typeHandler, validationHandler, attributeHandler, formHandler)
	r := setupMuxRouter(typeHandler, validationHandler, attributeHandler, formHandler)
	return r
}

func setupMuxRouter(typeHandler *handler.TypeHandler, validationHandler *handler.ValidationHandler,
	attributeHandler *handler.AttributeHandler, formHandler *handler.FormHandler) *mux.Router {

	api := mux.NewRouter()
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

	return api
}

func setupGinRouter(typeHandler *handler.TypeHandler, validationHandler *handler.ValidationHandler,
	attributeHandler *handler.AttributeHandler, formHandler *handler.FormHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/types/:id", gin.WrapF(typeHandler.GetType))
	r.POST("/types", gin.WrapF(typeHandler.CreateType))
	r.PUT("/types/:id", gin.WrapF(typeHandler.UpdateType))
	r.DELETE("/types/:id", gin.WrapF(typeHandler.DeleteType))

	r.GET("/validations/:id", gin.WrapF(validationHandler.GetValidation))
	r.POST("/validations", gin.WrapF(validationHandler.CreateValidation))
	r.PUT("/validations/:id", gin.WrapF(validationHandler.UpdateValidation))
	r.DELETE("/validations/:id", gin.WrapF(validationHandler.DeleteValidation))

	r.GET("/attributes/:id", gin.WrapF(attributeHandler.GetAttribute))
	r.POST("/attributes", gin.WrapF(attributeHandler.CreateAttribute))
	r.PUT("/attributes/:id", gin.WrapF(attributeHandler.UpdateAttribute))
	r.DELETE("/attributes/:id", gin.WrapF(attributeHandler.DeleteAttribute))

	r.GET("/forms/:id", gin.WrapF(formHandler.GetForm))
	r.POST("/forms", gin.WrapF(formHandler.CreateForm))
	r.PUT("/forms/:id", gin.WrapF(formHandler.UpdateForm))
	r.DELETE("/forms/:id", gin.WrapF(formHandler.DeleteForm))
	return r
}

func setupTestDB(logger *slog.Logger) *gorm.DB {
	dsn := "host=localhost user=test_public_config_user password=testpassword dbname=test_public_config_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect to database")
		panic(err)
	}
	db.AutoMigrate(&model.Type{}, &model.Validation{}, &model.Attribute{}, &model.Form{})
	return db
}

func TestTypeAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)
	createdType := model.Type{}

	t.Run("CreateType", func(t *testing.T) {
		newType := model.Type{
			Namespace:   "test_namespace",
			Family:      "test_family",
			Name:        "test_name",
			ElementType: "test_element",
			WidgetType:  "test_widget",
		}
		jsonValue, _ := json.Marshal(newType)
		req, _ := http.NewRequest("POST", "/types", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status code %d but got %d", http.StatusCreated, w.Code)
		} else {
			json.NewDecoder(w.Body).Decode(&createdType)
			t.Logf("created type %v", createdType)
		}
	})

	t.Run("GetType", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/types/%d", createdType.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotType model.Type
		json.Unmarshal(w.Body.Bytes(), &gotType)
		if gotType.ID != createdType.ID {
			t.Fatalf("expected ID %d but got %d", createdType.ID, gotType.ID)
		}
	})
}

func TestValidationAPI(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)
	createdValidation := model.Validation{}

	t.Run("CreateValidation", func(t *testing.T) {
		newValidation := model.Validation{
			Namespace:        "test_namespace",
			Family:           "test_family",
			Name:             "test_name",
			RuleName:         "test_rule",
			ValidationParams: "{}",
		}
		jsonValue, _ := json.Marshal(newValidation)
		req, _ := http.NewRequest("POST", "/validations", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status code %d but got %d", http.StatusCreated, w.Code)
		} else {
			json.NewDecoder(w.Body).Decode(&createdValidation)
			t.Logf("created validation %v", createdValidation)
		}
	})

	t.Run("GetValidation", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/validations/%d", createdValidation.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotValidation model.Validation
		json.Unmarshal(w.Body.Bytes(), &gotValidation)
		if gotValidation.ID != createdValidation.ID {
			t.Fatalf("expected ID %d but got %d", createdValidation.ID, gotValidation.ID)
		}
	})
}

func TestAttributeAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)
	createdAttribute := model.Attribute{}

	t.Run("CreateAttribute", func(t *testing.T) {
		newAttribute := model.Attribute{
			Namespace:  "test_namespace",
			Family:     "test_family",
			Name:       "test_name",
			Label:      "test_label",
			DesignSpec: "{\"color\":\"blue\"}",
		}
		jsonValue, _ := json.Marshal(newAttribute)
		req, _ := http.NewRequest("POST", "/attributes", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status code %d but got %d", http.StatusCreated, w.Code)
		} else {
			json.NewDecoder(w.Body).Decode(&createdAttribute)
			t.Logf("created attribute %v", createdAttribute)
		}
	})

	t.Run("GetAttribute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/attributes/%v", createdAttribute.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotAttribute model.Attribute
		json.Unmarshal(w.Body.Bytes(), &gotAttribute)
		if gotAttribute.ID != createdAttribute.ID {
			t.Fatalf("expected ID %d but got %d", createdAttribute.ID, gotAttribute.ID)
		}
	})
}

func TestFormAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)
	createdForm := model.Form{}

	t.Run("CreateForm", func(t *testing.T) {
		newForm := model.Form{
			Namespace:  "test_namespace",
			Family:     "test_family",
			Name:       "test_name",
			ActionName: "test_action",
			Attributes: "{}",
		}
		jsonValue, _ := json.Marshal(newForm)
		req, _ := http.NewRequest("POST", "/forms", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status code %d but got %d", http.StatusCreated, w.Code)
		} else {
			json.NewDecoder(w.Body).Decode(&createdForm)
			t.Logf("created form %v", createdForm)
		}
	})

	t.Run("GetForm", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/forms/%d", createdForm.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotForm model.Form
		json.Unmarshal(w.Body.Bytes(), &gotForm)
		if gotForm.ID != createdForm.ID {
			t.Fatalf("expected ID %d but got %d", createdForm.ID, gotForm.ID)
		}
	})
}
