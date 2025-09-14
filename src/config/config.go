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
		viper.SetEnvPrefix("SOCIALPLATFORM")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
