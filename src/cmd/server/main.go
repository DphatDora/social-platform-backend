package main

import (
	"os"
	"social-platform-backend/config"
	"social-platform-backend/internal/infrastructure/cache"
	"social-platform-backend/internal/infrastructure/db"
	"social-platform-backend/internal/interface/router"
	"social-platform-backend/internal/wire"
	"social-platform-backend/package/logger"
	"strconv"
	"time"
)

const (
	DefaultPort = 8045
)

func main() {
	// Set timezone to UTC
	time.Local = time.UTC

	conf := config.GetConfig()
	if err := logger.Init(&conf.Log); err != nil {
		logger.Errorf("[ERROR] Logger initialization failed: %v", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()
	logger.Debugf("[DEBUG] Config: %+v", conf)

	// init database
	db.InitPostgresql(&conf)
	defer func() {
		if err := db.ClosePostgresql(); err != nil {
			logger.Errorf("[ERROR] Close postgresql fail: %s\n", err)
		}
	}()

	// init Redis
	redisClient, err := cache.NewRedisClient(&conf)
	if err != nil {
		if conf.Redis.Required {
			logger.Fatalf("[ERROR] Redis is required but failed to initialize: %v", err)
		}
		logger.Warnf("[WARNING] Failed to initialize Redis: %v. Rate limiting and password cache disabled.", err)
		redisClient = nil
	}
	defer func() {
		if err := cache.CloseRedis(); err != nil {
			logger.Errorf("[ERROR] Close Redis fail: %s\n", err)
		}
	}()

	// wire-generated DI container
	appHandler := wire.InitAppContainer(db.GetDB(), redisClient, &conf)

	// set up routes
	r := router.SetupRoutes(appHandler, &conf)

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	logger.Infof("[✅] Server starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}
