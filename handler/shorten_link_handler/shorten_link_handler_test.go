package shortenlinkhandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mockService "github.com/huyvo1112/golang_course/mocks/shorten_link_service"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenLink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	serviceErr := errors.New("service error")
	// ctx := context.Background()

	testcases := []struct {
		name                 string
		requestBody          string
		setupMockService     func(t *testing.T, c *gin.Context) *mockService.MockShortenLink
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Normal case",
			requestBody: `{
		        "exp": 604800,
		        "url": "https://google.com"
		    }`,
			setupMockService: func(t *testing.T, c *gin.Context) *mockService.MockShortenLink {
				mockSvc := new(mockService.MockShortenLink)
				mockSvc.EXPECT().
					CreateShortLink(c, "https://google.com", 7*24*time.Hour).
					Return("abcd123", nil).
					Once()
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
		        "code": "abcd123",
		        "message": "Shorten URL generated successfully!"
		    }`,
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"exp": }`,

			setupMockService: func(_ *testing.T, _ *gin.Context) *mockService.MockShortenLink {
				return new(mockService.MockShortenLink)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Missing URL",

			requestBody: `{"exp": 604800}`,

			setupMockService: func(_ *testing.T, _ *gin.Context) *mockService.MockShortenLink {
				return new(mockService.MockShortenLink)
			},

			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid URL",
			requestBody: `{
		        "exp": 604800,
		        "url": "invalid-url"
		    }`,

			setupMockService: func(_ *testing.T, _ *gin.Context) *mockService.MockShortenLink {
				return new(mockService.MockShortenLink)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Service error",

			requestBody: `{
		        "exp": 604800,
		        "url": "https://example.com"
		    }`,

			setupMockService: func(t *testing.T, c *gin.Context) *mockService.MockShortenLink {
				mockSvc := new(mockService.MockShortenLink)

				mockSvc.EXPECT().
					CreateShortLink(
						c,
						"https://example.com",
						7*24*time.Hour,
					).
					Return("", serviceErr).
					Once()

				return mockSvc
			},

			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "EXP nil",

			requestBody: `{
		        "exp": null,
		        "url": "https://example.com"
		    }`,

			setupMockService: func(t *testing.T, c *gin.Context) *mockService.MockShortenLink {
				mockSvc := new(mockService.MockShortenLink)

				mockSvc.EXPECT().
					CreateShortLink(
						c,
						"https://example.com",
						time.Duration(0),
					).
					Return("abc1234", nil).
					Once()

				return mockSvc
			},

			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			t.Parallel()
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBufferString(tc.requestBody))
			c.Request = req

			mockSvc := tc.setupMockService(t, c)
			handler := NewShortenLink(mockSvc)

			handler.CreateShortenLink(c)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(
					t,
					tc.expectedResponseBody,
					rec.Body.String(),
				)
			}
			assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRedirectShortenLink(t *testing.T) {
	t.Parallel()
	serviceErr := errors.New("service error")
	testcases := []struct {
		name             string
		setupMockService func(ctx context.Context) *mockService.MockShortenLink
		setupTestRequest func(c *gin.Context)
		expectedCode     int
		expectUrl        string
	}{
		{
			name: "Normal case",
			setupMockService: func(ctx context.Context) *mockService.MockShortenLink {
				mockSvc := new(mockService.MockShortenLink)
				mockSvc.On("GetURLFromCode", ctx, "123456").Return("https://google.com", nil)
				return mockSvc
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/123456", nil)
				c.Params = gin.Params{{Key: "code", Value: "123456"}}
			},
			expectedCode: http.StatusMovedPermanently,
			expectUrl:    "https://google.com",
		},
		{
			name: "Missing code",
			setupMockService: func(_ context.Context) *mockService.MockShortenLink {
				return new(mockService.MockShortenLink)
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				c.Params = gin.Params{{Key: "code", Value: ""}}
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			setupMockService: func(ctx context.Context) *mockService.MockShortenLink {
				mockSvc := new(mockService.MockShortenLink)
				mockSvc.EXPECT().GetURLFromCode(ctx, "123456").Return("", serviceErr).Once()
				return mockSvc
			},
			setupTestRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/123456", nil)
				c.Params = gin.Params{{Key: "code", Value: "123456"}}
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			tc.setupTestRequest(c)

			mockSvc := tc.setupMockService(c)
			handler := NewShortenLink(mockSvc)
			handler.RedirectShortenLink(c)
			mockSvc.AssertExpectations(t)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectUrl, rec.Header().Get("Location"))
		})
	}
}
