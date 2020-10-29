package tserver

import (
	"log"
	"tbotext/trouter"

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

	userChan map[int]chan (tgbotapi.Message)
}

func (t *TelegramServer) getUserChan(ID int) chan (tgbotapi.Message) {
	if _, ok := t.userChan[ID]; !ok {
		t.userChan[ID] = make(chan tgbotapi.Message)
	}
	return t.userChan[ID]
}

// func (t *TelegramServer) getUserCtx(ID int) *context.Context {
// 	if _, ok := t.userCtx[ID]; !ok {
// 		log.Println("New user context")
// 		ctx := context.Background()
// 		t.userCtx[ID] = &ctx
// 	}
// 	log.Println("Old user context")
// 	return t.userCtx[ID]
// }

// ListenAndServe listens for incoming messages and serves them
func (t TelegramServer) ListenAndServe() {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	t.userChan = make(map[int]chan tgbotapi.Message)

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
		f := t.Router.Match(*update.Message)
		if f != nil {
			msgChan := t.getUserChan(update.Message.From.ID)
			go f.Serve(bot, msgChan)
			msgChan <- *update.Message
		}
	}
}
