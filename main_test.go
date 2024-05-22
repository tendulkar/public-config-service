package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"stellarsky.ai/platform/public-config-service/handler"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
	"stellarsky.ai/platform/public-config-service/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func convertToGinFunc(f func(http.ResponseWriter, *http.Request)) gin.HandlerFunc {

// 	return func(c *gin.Context) {
// 		f(c.Writer, c.Request)
// 	}
// }

func setupRouter(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

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
	r.GET("/types/:id", gin.WrapF(typeHandler.GetType))
	r.POST("/types", gin.WrapF(typeHandler.CreateType))
	r.PUT("/types/:id", gin.WrapF(typeHandler.UpdateType))
	r.DELETE("/types/:id", gin.WrapF(typeHandler.DeleteType))

	// Validation Routes
	r.GET("/validations/:id", gin.WrapF(validationHandler.GetValidation))
	r.POST("/validations", gin.WrapF(validationHandler.CreateValidation))
	r.PUT("/validations/:id", gin.WrapF(validationHandler.UpdateValidation))
	r.DELETE("/validations/:id", gin.WrapF(validationHandler.DeleteValidation))

	// Attribute Routes
	r.GET("/attributes/:id", gin.WrapF(attributeHandler.GetAttribute))
	r.POST("/attributes", gin.WrapF(attributeHandler.CreateAttribute))
	r.PUT("/attributes/:id", gin.WrapF(attributeHandler.UpdateAttribute))
	r.DELETE("/attributes/:id", gin.WrapF(attributeHandler.DeleteAttribute))

	// Form Routes
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
		panic("failed to connect to database")
	}
	db.AutoMigrate(&model.Type{}, &model.Validation{}, &model.Attribute{}, &model.Form{})
	return db
}

func TestTypeAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)

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

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetType", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/types/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotType model.Type
		json.Unmarshal(w.Body.Bytes(), &gotType)
		if gotType.ID != 1 {
			t.Fatalf("expected ID %d but got %d", 1, gotType.ID)
		}
	})
}

func TestValidationAPI(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)

	t.Run("CreateValidation", func(t *testing.T) {
		newValidation := model.Validation{
			Namespace:        "test_namespace",
			Family:           "test_family",
			Name:             "test_name",
			RuleName:         "test_rule",
			ValidationParams: "test_params",
		}
		jsonValue, _ := json.Marshal(newValidation)
		req, _ := http.NewRequest("POST", "/validations", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetValidation", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/validations/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotValidation model.Validation
		json.Unmarshal(w.Body.Bytes(), &gotValidation)
		if gotValidation.ID != 1 {
			t.Fatalf("expected ID %d but got %d", 1, gotValidation.ID)
		}
	})
}

func TestAttributeAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)

	t.Run("CreateAttribute", func(t *testing.T) {
		newAttribute := model.Attribute{
			Namespace:  "test_namespace",
			Family:     "test_family",
			Name:       "test_name",
			Label:      "test_label",
			DesignSpec: "test_spec",
		}
		jsonValue, _ := json.Marshal(newAttribute)
		req, _ := http.NewRequest("POST", "/attributes", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetAttribute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/attributes/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotAttribute model.Attribute
		json.Unmarshal(w.Body.Bytes(), &gotAttribute)
		if gotAttribute.ID != 1 {
			t.Fatalf("expected ID %d but got %d", 1, gotAttribute.ID)
		}
	})
}

func TestFormAPI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupTestDB(logger)
	router := setupRouter(db, logger)

	t.Run("CreateForm", func(t *testing.T) {
		newForm := model.Form{
			Namespace:  "test_namespace",
			Family:     "test_family",
			Name:       "test_name",
			ActionName: "test_action",
		}
		jsonValue, _ := json.Marshal(newForm)
		req, _ := http.NewRequest("POST", "/forms", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetForm", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/forms/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var gotForm model.Form
		json.Unmarshal(w.Body.Bytes(), &gotForm)
		if gotForm.ID != 1 {
			t.Fatalf("expected ID %d but got %d", 1, gotForm.ID)
		}
	})
}
