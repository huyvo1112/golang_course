package main

import (
	"log"

	"github.com/huyvo1112/golang_course/api"
	"github.com/huyvo1112/golang_course/config"
	"github.com/huyvo1112/golang_course/healthcheck_handler"
	"github.com/huyvo1112/golang_course/healthcheck_service"
)

func main() {
	// 1. Load environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize service
	service := healthcheck_service.NewHealthCheckService(
		cfg.ServiceName,
		cfg.InstanceID,
	)

	// 3. Initialize handler
	handler := healthcheck_handler.NewHealthCheckHandler(service)

	// 4. Initialize engine
	engine := api.NewEngine(handler, cfg.Port)

	// 5. Start server
	if err := engine.Start(); err != nil {
		log.Fatal(err)
	}
}
