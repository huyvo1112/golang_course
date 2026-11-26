package shortenlinkservice

import (
	"errors"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mocks "github.com/huyvo1112/golang_course/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestCreateShortLink(t *testing.T) {
	generatorErr := errors.New("generator error")
	repositoryErr := errors.New("repository error")
	testcases := []struct {
		name                   string
		url                    string
		expirationTime         time.Duration
		setupMockCodeGenerator func(t *testing.T) *mocks.MockCodeGenerator
		setupMockRepository    func(t *testing.T) *mocks.MockURLStorage
		expectedCode           string
		expectedError          error
	}{
		{
			name:           "Normal case",
			url:            "http://google.com",
			expirationTime: 7 * 24 * time.Hour,
			setupMockCodeGenerator: func(t *testing.T) *mocks.MockCodeGenerator {
				mockGenCode := new(mocks.MockCodeGenerator)
				mockGenCode.EXPECT().
					GenerateCode(7).
					Return("abcd123", nil).Once()
				return mockGenCode
			},
			setupMockRepository: func(t *testing.T) *mocks.MockURLStorage {
				mockRepo := new(mocks.MockURLStorage)
				mockRepo.EXPECT().
					SaveURL(mock.Anything,
						"abcd123",
						"http://google.com",
						7*24*time.Hour).
					Return(nil).Once()
				return mockRepo
			},
			expectedCode:  "abcd123",
			expectedError: nil,
		},
		{
			name:           "failed to generate short code: generator error",
			url:            "http://google.com",
			expirationTime: 7 * 24 * time.Hour,

			setupMockCodeGenerator: func(t *testing.T) *mocks.MockCodeGenerator {
				mockGenCode := new(mocks.MockCodeGenerator)

				mockGenCode.EXPECT().
					GenerateCode(7).
					Return("", generatorErr).
					Once()

				return mockGenCode
			},

			setupMockRepository: func(t *testing.T) *mocks.MockURLStorage {
				return new(mocks.MockURLStorage)
			},

			expectedCode:  "",
			expectedError: generatorErr,
		},
		{
			name:           "Repository error",
			url:            "http://google.com",
			expirationTime: 7 * 24 * time.Hour,

			setupMockCodeGenerator: func(t *testing.T) *mocks.MockCodeGenerator {
				mockGenCode := new(mocks.MockCodeGenerator)

				mockGenCode.EXPECT().
					GenerateCode(7).
					Return("abcd123", nil).
					Once()

				return mockGenCode
			},

			setupMockRepository: func(t *testing.T) *mocks.MockURLStorage {
				mockRepo := new(mocks.MockURLStorage)

				mockRepo.EXPECT().
					SaveURL(
						mock.Anything,
						"abcd123",
						"http://google.com",
						7*24*time.Hour,
					).
					Return(repositoryErr).
					Once()

				return mockRepo
			},

			expectedCode:  "",
			expectedError: repositoryErr,
		},
	}
	for _, tc := range testcases {
		gin.SetMode(gin.TestMode)
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGenCode := tc.setupMockCodeGenerator(t)
			mockRepo := tc.setupMockRepository(t)
			ctx := context.Background()
			svc := NewShortenLinkService(mockRepo, mockGenCode)
			actualCode, actualError := svc.CreateShortLink(ctx, tc.url, 7*24*time.Hour)

			assert.Equal(t, tc.expectedCode, actualCode)

			if tc.expectedError != nil {
				assert.ErrorIs(t, actualError, tc.expectedError)
			} else {
				assert.NoError(t, actualError)
			}

			// Verify mock expectations
			mockGenCode.AssertExpectations(t)
			mockRepo.AssertExpectations(t)

			// Generator error: Repository must not be called
			if tc.name == "Generator error" {
				mockRepo.AssertNotCalled(
					t,
					"SaveURL",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				)
			}
		})
	}
}
