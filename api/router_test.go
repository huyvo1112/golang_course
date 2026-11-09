package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/huyvo1112/golang_course/healthcheck_handler"
	"github.com/huyvo1112/golang_course/healthcheck_service"
	"github.com/huyvo1112/golang_course/model"
)

func TestHealthCheckRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Initialize service
	service := healthcheck_service.NewHealthCheckService(
		"bookmark_service",
		"123e4567-e89b-12d3-a456-426614174000",
	)

	// 2. Initialize handler
	handler := healthcheck_handler.NewHealthCheckHandler(service)

	// 3. Initialize router
	engine := NewEngine(handler, 8080)

	// 4. Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	recorder := httptest.NewRecorder()

	// 5. Send request to router
	engine.ServeHTTP(recorder, req)

	// 6. Assert HTTP status
	require.Equal(t, http.StatusOK, recorder.Code)

	// 7. Assert Content-Type
	assert.Equal(
		t,
		"application/json; charset=utf-8",
		recorder.Header().Get("Content-Type"),
	)

	// 8. Assert response body
	var actualResponse model.HealthCheck
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	require.NoError(t, err)

	expectedResponse := model.HealthCheck{
		Message:    "OK",
		Service:    "bookmark_service",
		InstanceID: "123e4567-e89b-12d3-a456-426614174000",
	}

	assert.Equal(t, expectedResponse, actualResponse)
}
