package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file\n v%", err.Error())
	}
}

const (
	CALENDAR string = "FRONTEND 3 - lunes, miercoles, jueves 18hsARG / 16hsCO"
	LINK string = "https://digitalhouse.zoom.us/my/aulavirtual94"
	start string = "BIENVENIDE\n\n" + help
	help string = "Comandos de ayuda\n" +
	"help - muestra la lista de comandos\n" +
	"open - abre el menu de opciones/teclado rapido\n" +
	"close - cierra el menu de opciones/teclado rapido\n" +
	"link - te dara el link de los encuentros de zoom\n" +
	"calendar - te mostrara las fechas importantes de la materia\n"
)

var count int8 = 3
var textResp string

func main() {
	runBot()
}

func runBot() {
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

		textResp = armarRespuesta(&update.Message.Text)

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
	switch strings.ToLower( *mensajeRecibido ){
	case "/start":
		respuesta = start
		count = 3
	case "help":
		respuesta = help
		count = 3
	case "calendar":
		respuesta = getEvents(GetOauthAndCalendar())
		count = 3
	case "link":
		respuesta = fmt.Sprintf("link:\n%v\npass:\n%v", LINK, os.Getenv("LINK_PASS"))
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
		tgbotapi.NewKeyboardButton("calendar"),
		tgbotapi.NewKeyboardButton("close"),
	),
)

func GetOauthAndCalendar() (*calendar.Service){
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithCredentialsFile("bbtl-api-conection.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return service
}

func getEvents(service *calendar.Service) string {
	t := time.Now().Format(time.RFC3339)
	events, err := service.Events.List(os.Getenv("calendarId")).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	var info string
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
		info = "No upcoming events found."
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			d, err := time.Parse(time.RFC3339, date)
			if err != nil {
				log.Fatalf("Error al parsear el tiempo: %v", err)
			}
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, d.Format(time.RFC850))
			info = strings.Join([]string{info, item.Summary, d.Format(time.RFC850), "\n"}, " ")
		}
	}
	return info
}
