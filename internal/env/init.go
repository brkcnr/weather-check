package env

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("cmd/server/.env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}
}
