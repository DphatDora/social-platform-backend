package config

import (
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	App      App
	Database Database
	Log      Log
	Auth     Auth
	Client   Client
	Server   Server
	Redis    Redis
}

func LoadConfig() {
	once.Do(func() {
		// load .env config
		_ = godotenv.Load()

		// load /yaml config
		viper.SetConfigName("config") // "dev"
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Config file error: %s", err)
		}

		// bind system environment variables
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		bindEnvs()

		// load into struct
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Config unmarshal error: %s", err)
		}
	})
}

func GetConfig() Config {
	LoadConfig()
	return config
}

func bindEnvs() {
	// App
	_ = viper.BindEnv("app.name", "APP_NAME")
	_ = viper.BindEnv("app.host", "HOST")
	_ = viper.BindEnv("app.port", "PORT")

	// Database
	_ = viper.BindEnv("database.username", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.name", "DB_NAME")
	_ = viper.BindEnv("database.sslMode", "DB_SSLMODE")
	_ = viper.BindEnv("database.timeZone", "DB_TIMEZONE")

	// Server
	_ = viper.BindEnv("server.url", "SERVER_URL")

	// Client
	_ = viper.BindEnv("client.url", "CLIENT_URL")

	// Auth
	_ = viper.BindEnv("auth.jwtSecret", "JWT_SECRET")
	_ = viper.BindEnv("auth.googleClientID", "GOOGLE_CLIENT_ID")
	_ = viper.BindEnv("auth.googleClientSecret", "GOOGLE_CLIENT_SECRET")
	_ = viper.BindEnv("auth.googleRedirectURL", "GOOGLE_REDIRECT_URL")

	// Redis
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "REDIS_DB")
	_ = viper.BindEnv("redis.poolSize", "REDIS_POOL_SIZE")
	_ = viper.BindEnv("redis.required", "REDIS_REQUIRED")
}
