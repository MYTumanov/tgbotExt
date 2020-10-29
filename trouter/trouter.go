package trouter

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TbotServe ...
type TbotServe func(*tgbotapi.BotAPI, tgbotapi.Message)

// TbotChanServe ...
type TbotChanServe func(*tgbotapi.BotAPI, <-chan (tgbotapi.Message))

type Handler interface {
	Serve(*tgbotapi.BotAPI, <-chan (tgbotapi.Message))
}

// Serve ...
func (t TbotServe) Serve(b *tgbotapi.BotAPI, mChan <-chan (tgbotapi.Message)) {
	m := <-mChan
	t(b, m)
}

// Serve ...
func (t TbotChanServe) Serve(b *tgbotapi.BotAPI, mChan <-chan (tgbotapi.Message)) {

	t(b, mChan)
}

// Route ...
type Route struct {
	command         string
	commandFunc     Handler
	chain           []Handler
	curChainElement int
}

// Router ...
type Router struct {
	// Stores commands and func to handle
	// router map[string]TbotServe

	// Stores user id and func that must be handle to request message
	chainHandlers map[int]*Route

	// Stores commands and chained commands
	// chain map[string]Route

	// Stores commands and hadle struct
	handlers map[string]*Route
}

// NewRouter ...
func NewRouter() *Router {
	return &Router{
		// router:      make(map[string]TbotServe),
		chainHandlers: make(map[int]*Route),
		// chain:       make(map[string]Route),
		handlers: make(map[string]*Route),
	}
}

// Match ...
func (r *Router) Match(msg tgbotapi.Message) Handler {
	log.Printf("MATCH: start, text %v \n", msg.Text)
	if r.handlers == nil {
		return nil
	}
	userID := msg.From.ID

	if msg.IsCommand() {
		log.Printf("MATCH: command true, text %v \n", msg.Text)
		if route, ok := r.handlers[msg.Text]; ok {
			log.Printf("MATCH: matched true \n")
			r.chainHandlers[userID] = route
			return route.commandFunc
		}
		log.Printf("MATCH: matched false \n")
	} else {
		log.Printf("MATCH: command false, text %v \n", msg.Text)
		if route, ok := r.chainHandlers[userID]; ok {
			log.Printf("MATCH: user matched true \n")
			var f Handler
			if len(route.chain) > route.curChainElement {
				log.Printf("MATCH: chained func found true \n")
				f = route.chain[route.curChainElement]
				route.curChainElement++
			}
			if len(route.chain) == route.curChainElement {
				delete(r.chainHandlers, userID)
			}
			return f
		}
		log.Printf("MATCH: user matched false \n")
	}
	return nil
}

// HandleComandFunc adds commad and handle func
func (r *Router) HandleComandFunc(command string, f func(*tgbotapi.BotAPI, tgbotapi.Message)) *Route {
	route := &Route{
		command:         command,
		commandFunc:     TbotServe(f),
		curChainElement: 0,
	}
	// r.chain[command] = *route
	r.handlers[command] = route
	return route
}

// ChainedFunc adds chain handle func
func (r *Route) ChainedFunc(f func(*tgbotapi.BotAPI, tgbotapi.Message)) *Route {
	r.chain = append(r.chain, TbotServe(f))
	return r
}
