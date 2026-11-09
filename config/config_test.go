package config

import (
	"testing"

	"github.com/google/uuid"
)

func TestLoadConfig(t *testing.T) {

	testcases := []struct {
		name              string
		serviceName       string
		instanceID        string
		port              string
		expectError       bool
		expectServiceName string
		expectInstanceID  string
		expectPort        int
	}{
		{
			name:              "default config",
			serviceName:       "",
			instanceID:        "",
			port:              "",
			expectError:       false,
			expectServiceName: "bookmark_service",
			expectInstanceID:  "",
			expectPort:        8080,
		},
		{
			name:              "custom config",
			serviceName:       "my_service",
			instanceID:        "123e4567-e89b-12d3-a456-426614174000",
			port:              "9090",
			expectError:       false,
			expectServiceName: "my_service",
			expectInstanceID:  "123e4567-e89b-12d3-a456-426614174000",
			expectPort:        9090,
		},
		{
			name:        "invalid instance ID",
			serviceName: "my_service",
			instanceID:  "invalid-id",
			port:        "9090",
			expectError: true,
		},
		{
			name:        "invalid port",
			serviceName: "my_service",
			instanceID:  "123e4567-e89b-12d3-a456-426614174000",
			port:        "invalid-port",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for the test case
			if tc.serviceName != "" {
				t.Setenv("SERVICE_NAME", tc.serviceName)
			}

			if tc.instanceID != "" {
				t.Setenv("INSTANCE_ID", tc.instanceID)
			}

			if tc.port != "" {
				t.Setenv("PORT", tc.port)
			}

			// Call the function being tested
			config, err := LoadConfig()

			// Check the results
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				// Check InstanceID
				if tc.instanceID != "" {
					// Custom InstanceID: should match the input
					if config.InstanceID != tc.instanceID {
						t.Errorf("Expected instance ID %s, got %s", tc.instanceID, config.InstanceID)
					}
				} else {
					// Default InstanceID: should be a valid generated UUID
					if _, err := uuid.Parse(config.InstanceID); err != nil {
						t.Errorf("Expected a valid generated UUID, got %s", config.InstanceID)
					}
				}
				if config.ServiceName != tc.expectServiceName {
					t.Errorf("Expected service name %s, but got %s", tc.expectServiceName, config.ServiceName)
				}

				if config.Port != tc.expectPort {
					t.Errorf("Expected port %d, but got %d", tc.expectPort, config.Port)
				}
			}
		})
	}
}
