package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var textResp = armarRespuesta(&update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, textResp)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func start() string {
	return "bienvenido"
}

func showConfig() string {
	return "mostrando configuracion"
}

func armarRespuesta(mensajeRecibido *string) string {
	var respuesta string
	switch *mensajeRecibido {
	case "/start":
		respuesta = start()
	case "/config":
		respuesta = showConfig()
	case "ctd":
		respuesta = "FRONTEND 3 - lunes, miercoles, jueves 18hsARG / 16hsCO"
	default:
		respuesta = "respuesta generica"
	}
	return respuesta
}
