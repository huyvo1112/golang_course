package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huyvo1112/golang_course/generator"
	"github.com/huyvo1112/golang_course/handler/healthcheck_handler"
	shortenlinkhandler "github.com/huyvo1112/golang_course/handler/shorten_link_handler"
	"github.com/huyvo1112/golang_course/repository"
	shortenlinkservice "github.com/huyvo1112/golang_course/service/shorten_link_service"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app         *gin.Engine
	port        int
	redisClient *redis.Client
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func NewEngine(handler healthcheck_handler.HealthCheckHandler, port int, redisClient *redis.Client) Engine {
	e := &engine{
		app:         gin.Default(),
		port:        port,
		redisClient: redisClient,
	}

	e.initializeRoutes(handler)
	return e
}

func (e *engine) Start() error {
	address := ":" + strconv.Itoa(e.port)
	return e.app.Run(address)
}

func (e *engine) initializeRoutes(handler healthcheck_handler.HealthCheckHandler) {
	genCode := generator.NewCodeGenerator()
	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	shortenURLService := shortenlinkservice.NewShortenLinkService(urlStorageRepo, genCode)
	shortenLinkHandler := shortenlinkhandler.NewShortenLink(shortenURLService)

	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", handler.HandleHealthCheck)
	e.app.POST("v1/links/shorten", shortenLinkHandler.CreateShortenLink)
	e.app.GET("/v1/links/redirect/:code", shortenLinkHandler.RedirectShortenLink)
}
