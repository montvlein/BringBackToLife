package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	log.Println(token)
}