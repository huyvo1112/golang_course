package healthcheck_service

import (
	"testing"

	"github.com/huyvo1112/golang_course/model"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthCheck(t *testing.T) {
	testcases := []struct {
		name           string
		message        string
		serviceName    string
		instanceID     string
		expectedValues *model.HealthCheck
		expectedErr    bool
	}{
		{
			name:        "Valid health check",
			message:     "OK",
			serviceName: "bookmark-service",
			instanceID:  "18ee44d6-3010-4587-a2b5-0a4de762a738",
			expectedValues: &model.HealthCheck{
				Message:    "OK",
				Service:    "bookmark-service",
				InstanceID: "18ee44d6-3010-4587-a2b5-0a4de762a738",
			},
			expectedErr: false,
		},
		{
			name:        "Empty service name",
			serviceName: "",
			instanceID:  "123e4567-e89b-12d3-a456-426614174000",
			expectedValues: &model.HealthCheck{
				Message:    "OK",
				Service:    "",
				InstanceID: "123e4567-e89b-12d3-a456-426614174000",
			},
			expectedErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := NewHealthCheckService(tc.serviceName, tc.instanceID)

			healthCheck, err := service.GetHealthCheck()

			if (err != nil) != tc.expectedErr {
				t.Fatalf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
			if healthCheck == nil {
				t.Fatalf("Expected health check, got nil")
			}
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if assert.NotNil(t, healthCheck) {
				assert.Equal(t, tc.expectedValues.Message, healthCheck.Message)
				assert.Equal(t, tc.expectedValues.Service, healthCheck.Service)
				assert.Equal(t, tc.expectedValues.InstanceID, healthCheck.InstanceID)
			}
		})
	}
}
