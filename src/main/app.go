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
		log.Fatalf("Error loading .env file\n %v", err.Error())
	}
}

const (
	start string = "BIENVENIDE\n\n" + help
	help string = "Comandos de ayuda\n" +
	"help - muestra la lista de comandos\n" +
	"open - abre el menu de opciones/teclado rapido\n" +
	"close - cierra el menu de opciones/teclado rapido\n" +
	"siguiente - te mostrara el evento mas proximo\n" +
	"proximos - te mostrara las fechas importantes de la materia\n" +
	"agregar - crea un evento nuevo\n"
)

var count int8 = 3
var textResp string
var calendarApi googleCalendar = googleCalendar{service: GetOauthAndCalendar()}

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

	defer avisarProximoEvento(bot) // al entrar al for de updates no entra a esta funcion. Se necesita correr ambas en paralelo

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
	case "siguiente":
		respuesta = calendarApi.getNextEvent()
		count = 3
	case "proximos":
		respuesta = calendarApi.getEvents()
		count = 3
	case "crear":
		calendarApi.createEvent(
			"Google I/O 2015",
			"800 Howard St., San Francisco, CA 94103",
			"A chance to hear more about Google's developer products.",
			"2015-05-28T09:00:00-07:00",
			"2015-05-28T17:00:00-07:00")
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

type googleCalendar struct{
	service *calendar.Service
}

func GetOauthAndCalendar() (*calendar.Service){
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithCredentialsFile("bbtl-api-conection.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return service
}

func (calApi googleCalendar) getEvents() string {
	t := time.Now().Format(time.RFC3339)
	events, err := calApi.service.Events.List(os.Getenv("calendarId")).ShowDeleted(false).
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

func (calApi googleCalendar) getNextEvent() string {
	t := time.Now().Format(time.RFC3339)
	events, err := calApi.service.Events.List(os.Getenv("calendarId")).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next user's event: %v", err.Error())
	}
	var info string
	if len(events.Items) == 0 {
		info = "No se encontraron eventos registrados."
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
			info = strings.Join([]string{info, item.Summary, "\n", d.Format(time.RFC850), "\n"}, " ")
		}
	}
	return info
}

func (calApi googleCalendar) getTimeNextEvent() string {
	t := time.Now().Format(time.RFC3339)
	events, err := calApi.service.Events.List(os.Getenv("calendarId")).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next user's event: %v", err.Error())
	}
	var info string
	if len(events.Items) == 0 {
		info = "No se encontraron eventos registrados."
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			info=date
		}
	}
	return info
}

func (calApi googleCalendar) createEvent(titulo, lugar, descripcion, inicio, fin string ) {

	eventSummary := &calendar.Event{
	Summary: titulo, // "Google I/O 2015"
	Location: lugar, // "800 Howard St., San Francisco, CA 94103"
	Description: descripcion, // "A chance to hear more about Google's developer products."
	Start: &calendar.EventDateTime{
	  DateTime: inicio, // "2015-05-28T09:00:00-07:00"
	  TimeZone: "America/Los_Angeles",
	},
	End: &calendar.EventDateTime{
	  DateTime: fin, // "2015-05-28T17:00:00-07:00"
	  TimeZone: "America/Los_Angeles",
	},
	Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
	Attendees: []*calendar.EventAttendee{
	  &calendar.EventAttendee{Email:"lpage@example.com"},
	  &calendar.EventAttendee{Email:"sbrin@example.com"},
	},
  }

  event, err := calApi.service.Events.Insert(os.Getenv("calendarId"), eventSummary).Do()
  if err != nil {
	log.Fatalf("Unable to create event. %v\n", err)
  }
  fmt.Printf("Event created: %s\n", event.HtmlLink)

}

func avisarProximoEvento(bot *tgbotapi.BotAPI) {
	chat, _ := strconv.Atoi(os.Getenv("CHAT_ID"))

	for range time.Tick(60 * time.Second) {
		proximo := calendarApi.getNextEvent()
		fechaProximo := calendarApi.getTimeNextEvent()
		fecha, _ := time.Parse(time.RFC3339,fechaProximo)
		diferencia := time.Until(fecha)
		minutosFaltantes := int(diferencia.Minutes())
		aviso := 30
		log.Printf("faltan %v minutos para enviar el mensaje...\n", minutosFaltantes)
		if minutosFaltantes == aviso {
			log.Println("Faltan 30 min para el evento. Enviando mensaje")
			mensaje := tgbotapi.NewMessage(int64(chat),proximo)
			bot.Send(mensaje)
			aviso = 0
		}
	}
}