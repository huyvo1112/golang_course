package healthcheck_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huyvo1112/golang_course/healthcheck_service"
)

type HealthCheckHandler interface {
	HandleHealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	healthCheckService healthcheck_service.HealthCheckService
}

func NewHealthCheckHandler(healthCheckService healthcheck_service.HealthCheckService) HealthCheckHandler {
	return &healthCheckHandler{
		healthCheckService: healthCheckService,
	}
}

func (h *healthCheckHandler) HandleHealthCheck(c *gin.Context) {
	healthCheck, err := h.healthCheckService.GetHealthCheck()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, healthCheck)
}
