package main

import (
	"log"

	"github.com/huyvo1112/golang_course/api"
	"github.com/huyvo1112/golang_course/config"
	"github.com/huyvo1112/golang_course/handler/healthcheck_handler"
	redisclient "github.com/huyvo1112/golang_course/pkg/redis"
	"github.com/huyvo1112/golang_course/service/healthcheck_service"

	_ "github.com/huyvo1112/golang_course/docs"
)

// @title           Bookmark Management API
// @version         1.0.0
// @description     API for managing bookmark
// @BasePath        /

func main() {
	// 1. Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize redisClient
	redisClent, err := redisclient.NewClient("")
	if err != nil {
		log.Fatalf("failed to initialize Redis client: %v", err)
	}

	// 3. Initialize service
	service := healthcheck_service.NewHealthCheckService(
		cfg.ServiceName,
		cfg.InstanceID,
	)

	// 4. Initialize handler
	handler := healthcheck_handler.NewHealthCheckHandler(service)

	// 5. Initialize engine
	engine := api.NewEngine(handler, cfg.Port, redisClent)

	// 6. Start server
	if err := engine.Start(); err != nil {
		log.Fatal(err)
	}
}
