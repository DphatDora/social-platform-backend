package db

import (
	"fmt"
	"social-platform-backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgSingleton *gorm.DB

func InitPostgresql(conf *config.Config) {
	dbUser := conf.Database.Username
	dbPassword := conf.Database.Password
	dbHost := conf.Database.Host
	dbPort := conf.Database.Port
	dbName := conf.Database.Name
	dbSslMode := conf.Database.SslMode
	dbTimezone := conf.Database.TimeZone

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode, dbTimezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db = db.Debug()
	if err != nil {
		panic("❌❌ Failed to connect database" + err.Error())
	}
	fmt.Println("✅✅ Connect to the database successfully")
	pgSingleton = db
}

func GetDB() *gorm.DB {
	if pgSingleton == nil {
		panic("Connection to Database Postgres is not setup")
	}

	return pgSingleton
}

func ClosePostgresql() error {
	sqlDB, err := pgSingleton.DB()
	if err != nil {
		fmt.Println("failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	return err
}
