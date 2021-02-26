package trouter

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// HandlerFunc is a type of executed func
type HandlerFunc func(*tgbotapi.BotAPI, *tgbotapi.Message)

type Router struct {
	// handler keeps command and chain of funcs
	handlers map[string]*Handler

	// context keep user and its function/chainfunctions to run
	uContxt map[int]*userContext
}

// Handler keeps list of funcs
type Handler struct {
	funcs []HandlerFunc
}

type userContext struct {
	userID         int
	currentFunc    HandlerFunc
	currentFuncNum int
	currentCommand string

	router *Router
}

// NewRouter returns new router
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]*Handler),
		uContxt:  make(map[int]*userContext),
	}
}

// HandleComandFunc adds commad and handle func type func(*tgbotapi.BotAPI, string)
func (r *Router) HandleComandFunc(command string, f HandlerFunc) *Handler {
	var tmpHandler []HandlerFunc
	tmpHandler = append(tmpHandler, f)
	handler := Handler{
		funcs: tmpHandler,
	}
	r.handlers[command] = &handler
	return &handler
}

// ChainedFunc adds chain handle func
func (h *Handler) ChainedFunc(f HandlerFunc) *Handler {
	h.funcs = append(h.funcs, f)
	return h
}

// Match search for func by input command and run it
func (r *Router) Match(command string, userID int) (HandlerFunc, error) {
	// define is user context exist, if not - set it
	if _, ok := r.uContxt[userID]; !ok {
		r.uContxt[userID] = &userContext{
			userID:         userID,
			router:         r,
			currentCommand: command,
		}
	}
	userCnxt := r.uContxt[userID]

	if command != "" {
		command = "/" + command
		if _, ok := r.handlers[command]; !ok {
			return nil, errors.New("Command not found")
		}

		userCnxt.currentFuncNum = 0
		userCnxt.currentCommand = command
		userCnxt.currentFunc = r.handlers[command].funcs[userCnxt.currentFuncNum]

	} else {
		userCnxt.currentFuncNum++
		if len(r.handlers[userCnxt.currentCommand].funcs) > userCnxt.currentFuncNum {
			userCnxt.currentFunc = r.handlers[userCnxt.currentCommand].funcs[userCnxt.currentFuncNum]
		} else {
			userCnxt.currentFunc = nil
			return nil, errors.New("Command not found")
		}

	}

	return userCnxt.currentFunc, nil
}

func (u *userContext) getExecutionFunc() HandlerFunc {
	u.currentFuncNum++
	if len(u.router.handlers[u.currentCommand].funcs) > u.currentFuncNum {
		u.currentFunc = u.router.handlers[u.currentCommand].funcs[u.currentFuncNum]

	} else {
		u.currentFunc = nil
	}
	return u.currentFunc
}
