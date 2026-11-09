package healthcheck_service

import "github.com/huyvo1112/golang_course/model"

type HealthCheckService interface {
	GetHealthCheck() (*model.HealthCheck, error)
}

type healthCheckService struct {
	serviceName string
	instanceID  string
}

func NewHealthCheckService(serviceName, instanceID string) HealthCheckService {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (h *healthCheckService) GetHealthCheck() (*model.HealthCheck, error) {
	return &model.HealthCheck{
		Message:    "OK",
		Service:    h.serviceName,
		InstanceID: h.instanceID,
	}, nil
}
