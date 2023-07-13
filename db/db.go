package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func handleError(err error, str string) {
	if err != nil {
		log.Fatalf("error has happened %s %v", str, err)
	}
}

func initDatabase() *gorm.DB {
	var err error
	dsn := "host=localhost user=gorm password=gorm dbname=invoicer port=9920 sslmode=disable TimeZone=UTC"
	DBConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	handleError(err, "Connection Opened to Database")
	DBConn.AutoMigrate(&Invoice{})

	return DBConn
	// don't forget to defer dbconn.close() in main
}
