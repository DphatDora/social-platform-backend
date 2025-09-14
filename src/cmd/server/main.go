package main

import (
	"fmt"
	"log"
	"social-platform-backend/config"
	"social-platform-backend/internal/infrastructure/db"
	"social-platform-backend/internal/interface/router"
	"strconv"
)

const (
	DefaultPort = 8045
)

func main() {
	setUpInfrastructure()
	defer closeInfrastructure()
}

func setUpInfrastructure() {
	conf := config.GetConfig()
	fmt.Println("[DEBUG] Config:", conf)

	// init database
	db.InitPostgresql(&conf)

	// set up routes
	r := router.SetupRoutes(db.GetDB(), &conf)

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("[✅✅] Server starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}

func closeInfrastructure() {
	if err := db.ClosePostgresql(); err != nil {
		log.Printf("[ERROR] Close postgresql fail: %s\n", err)
	}
}
