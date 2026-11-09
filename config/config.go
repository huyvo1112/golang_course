package config

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

func generateUUID() string {
	return uuid.NewString()
}

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process configuration: %w", err)
	}

	instanceID := strings.TrimSpace(cfg.InstanceID)
	cfg.InstanceID = instanceID
	if instanceID != "" {
		// Validate the UUID format
		if _, err := uuid.Parse(instanceID); err != nil {
			return nil, fmt.Errorf("invalid INSTANCE_ID format: %w", err)
		}
	} else {
		cfg.InstanceID = generateUUID()
	}

	return cfg, nil
}
