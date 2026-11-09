package model

type HealthCheck struct {
	Message    string `json:"message"`
	Service    string `json:"service"`
	InstanceID string `json:"instance_id"`
}
