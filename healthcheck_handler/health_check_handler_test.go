package healthcheck_handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	mockService "github.com/huyvo1112/golang_course/mocks/healthcheck_service"
	"github.com/huyvo1112/golang_course/model"
)

func TestHealthCheckHandler(t *testing.T) {

	testcases := []struct {
		name                 string
		setupMockService     func(t *testing.T) *mockService.MockHealthCheckService
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Health check successful",
			setupMockService: func(t *testing.T) *mockService.MockHealthCheckService {
				mockSvc := mockService.NewMockHealthCheckService(t)
				mockSvc.EXPECT().
					GetHealthCheck().
					Return(&model.HealthCheck{
						Message:    "OK",
						Service:    "bookmark-service",
						InstanceID: "123e4567-e89b-12d3-a456-426614174000",
					}, nil).
					Once()
				return mockSvc
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK","service":"bookmark-service","instance_id":"123e4567-e89b-12d3-a456-426614174000"}`,
		},
	}

	for _, tc := range testcases {
		gin.SetMode(gin.TestMode)
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mocksv := tc.setupMockService(t)
			handler := NewHealthCheckHandler(mocksv)
			router := gin.Default()
			router.GET("/health-check", handler.HandleHealthCheck)
			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
		})
	}
}
