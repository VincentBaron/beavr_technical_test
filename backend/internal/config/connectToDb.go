package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var DB_URL string

func ConnectToDb() {
	var err error
	// host := Conf.Database.Host
	// port := Conf.Database.Port
	// user := Conf.Database.User
	// password := Conf.Database.Password
	// database := Conf.Database.Name
	flyURI := os.Getenv("DATABASE_URL")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Connect to postgres
	dsn := fmt.Sprintf("%s", flyURI)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
}
