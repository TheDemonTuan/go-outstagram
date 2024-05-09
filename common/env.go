package common

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvVar() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
