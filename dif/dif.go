package dif

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func getDatabaseHost() string {
	if os.Getenv("MODE") == "production" {
		return "aumsu_postgres"
	}
	return "localhost"
}

var connStr = "host=" + getDatabaseHost() + " user=postgres password=postgres dbname=postgres sslmode=disable"
var DB, DBError = gorm.Open(postgres.New(postgres.Config{
	DSN: connStr,
}), &gorm.Config{})
