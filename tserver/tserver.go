package tserver

import (
	"log"

	"github.com/MYTumanov/tgbotext/trouter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramServer is struct with parameters
type TelegramServer struct {
	Token string

	WebHook         string
	ListenToWebhook string

	// Timeout uses when webhook dones't set
	Timeout int

	Router trouter.Router
}

// ListenAndServe listens for incoming messages and serves them
func (t TelegramServer) ListenAndServe() {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	// webhook or going for update with timeout
	var updates tgbotapi.UpdatesChannel
	if t.WebHook != "" {
		_, err := bot.SetWebhook(tgbotapi.NewWebhook(t.WebHook))
		if err != nil {
			log.Panic(err)
		}

		updates = bot.ListenForWebhook(t.ListenToWebhook)
	} else {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = t.Timeout
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			log.Panic(err)
		}
	}

	for update := range updates {
		log.Println("LOG ", update.Message.From.ID)
		log.Println("LOG ", update.Message.Text)
		f, err := t.Router.Match(update.Message.Text, update.Message.From.ID)
		if err != nil {
			log.Println(err)
		} else {
			f(bot, update.Message)
		}
	}
}
