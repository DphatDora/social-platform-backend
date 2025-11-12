package main

import (
	"fmt"
	"log"
	"social-platform-backend/config"
	"social-platform-backend/internal/infrastructure/cache"
	"social-platform-backend/internal/infrastructure/db"
	"social-platform-backend/internal/interface/router"
	"strconv"
	"time"
)

const (
	DefaultPort = 8045
)

func main() {
	setUpInfrastructure()
	defer closeInfrastructure()
}

func setUpInfrastructure() {
	// Set timezone to UTC
	time.Local = time.UTC

	conf := config.GetConfig()
	fmt.Println("[DEBUG] Config:", conf)

	// init database
	db.InitPostgresql(&conf)

	// init Redis
	redisClient, err := cache.NewRedisClient(&conf)
	if err != nil {
		if conf.Redis.Required {
			log.Fatalf("[ERROR] Redis is required but failed to initialize: %v", err)
		}
		log.Printf("[WARNING] Failed to initialize Redis: %v. Rate limiting and password cache disabled.", err)
		redisClient = nil
	}

	// set up routes
	r := router.SetupRoutes(db.GetDB(), redisClient, &conf)

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("[✅✅] Server starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}

func closeInfrastructure() {
	if err := cache.CloseRedis(); err != nil {
		log.Printf("[ERROR] Close Redis fail: %s\n", err)
	}

	if err := db.ClosePostgresql(); err != nil {
		log.Printf("[ERROR] Close postgresql fail: %s\n", err)
	}
}
