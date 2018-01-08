package internal

type commandHandler struct {
}

var commandHandlers = map[string]commandHandler{
	"start-ci": {},
}
