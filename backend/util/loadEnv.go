package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
var (
	Mongo_uri string
	Jwt_Secret string
)
func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Mongo_uri=os.Getenv("MONGO_URI")
	Jwt_Secret=os.Getenv("JWT_SECRET")
}