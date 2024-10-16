package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"Evolution/internal/app/ds"
	"Evolution/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&ds.Cards{}, &ds.Users{}, &ds.Moves{}, &ds.MoveCards{})
	if err != nil {
		panic("cant migrate db")
	}
}
