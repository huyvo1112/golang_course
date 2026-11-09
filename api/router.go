package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huyvo1112/golang_course/healthcheck_handler"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app  *gin.Engine
	port int
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func NewEngine(handler healthcheck_handler.HealthCheckHandler, port int) Engine {
	e := &engine{
		app:  gin.Default(),
		port: port,
	}
	e.initializeRoutes(handler)
	return e
}

func (e *engine) Start() error {
	address := ":" + strconv.Itoa(e.port)
	return e.app.Run(address)
}

func (e *engine) initializeRoutes(handler healthcheck_handler.HealthCheckHandler) {
	e.app.GET("/health-check", handler.HandleHealthCheck)
}
