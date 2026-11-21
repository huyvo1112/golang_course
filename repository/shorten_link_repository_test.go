package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage(t *testing.T) {
	t.Parallel()
	redisErr := errors.New("redis error")

	testcases := []struct {
		name          string
		code          string
		url           string
		setupMock     func(redisMock redismock.ClientMock)
		expectedError error
		verifyFunc    func(ctx context.Context, redisMock redismock.ClientMock)
	}{
		{
			name: "Normal case",
			code: "abcxyz",
			url:  "https://google.com",
			setupMock: func(redisMock redismock.ClientMock) {
				redisMock.ExpectSet(
					"abcxyz",
					"https://google.com",
					7*24*time.Hour,
				).SetVal("OK")
			},
			expectedError: nil,
			verifyFunc: func(ctx context.Context, redisMock redismock.ClientMock) {
				assert.NoError(t, redisMock.ExpectationsWereMet())
			},
		},
		{
			name:          "Redis error",
			code:          "abcxyz",
			url:           "https://google.com",
			expectedError: redisErr,
			setupMock: func(redisMock redismock.ClientMock) {
				redisMock.ExpectSet(
					"abcxyz",
					"https://google.com",
					7*24*time.Hour,
				).SetErr(redisErr)
			},
			verifyFunc: func(ctx context.Context, redisMock redismock.ClientMock) {
				assert.NoError(t, redisMock.ExpectationsWereMet())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisClient, redisMock := redismock.NewClientMock()
			defer redisClient.Close()
			tc.setupMock(redisMock)

			testRepo := NewURLStorage(redisClient)

			// Act
			err := testRepo.SaveURL(ctx, tc.code, tc.url, 7*24*time.Hour)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}
