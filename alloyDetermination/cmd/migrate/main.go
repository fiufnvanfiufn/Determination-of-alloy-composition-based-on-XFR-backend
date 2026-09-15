package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"alloyDetermination/internal/app/ds"
	"alloyDetermination/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	fmt.Println("DSN:", dsn.FromEnv())
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Alloy{},
		&ds.Like{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
