package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		switch update.Message.Text {
		case "/start", "open":
			msg.ReplyMarkup = optionKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func armarRespuesta(mensajeRecibido *string) string {
	var respuesta string
	switch *mensajeRecibido {
	case "/start":
		respuesta = start
		count = 3
	case "help":
		respuesta = help
		count = 3
	case "calendar":
		respuesta = os.Getenv("CALENDAR")
		count = 3
	case "link":
		respuesta = fmt.Sprintf("link:\n%v\npass:\n%v", os.Getenv("LINK"), os.Getenv("LINK_PASS"))
		count = 3
	case "open":
		respuesta = "menu abierto"
	case "close":
		respuesta = "menu cerrado"

	default:
		if count > 0 {
			count--
			respuesta = "No entiendo la consulta... \nEste mensaje se mostrara hasta " + strconv.Itoa(int(count)) + " veces seguidas mas"
		}
	}
	return respuesta
}

var optionKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("help"),
		tgbotapi.NewKeyboardButton("link"),
		tgbotapi.NewKeyboardButton("calendar"),
	),
)

var start string = "BIENVENIDE\n\n" + help

var help string = "Comandos de ayuda\n" +
	"help - muestra la lista de comandos\n" +
	"open - abre el menu de opciones/teclado rapido\n" +
	"close - cierra el menu de opciones/teclado rapido\n" +
	"link - te dara el link de los encuentros de zoom\n" +
	"calendar - te mostrara las fechas importantes de la materia\n"

var count int8 = 3
