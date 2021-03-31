package dif

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connStr = "host=188.225.47.219 port=6432 user=postgres password=postgres dbname=postgres sslmode=disable"
var DB, DBError = gorm.Open(postgres.New(postgres.Config{
	DSN: connStr,
}), &gorm.Config{})
