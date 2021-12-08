package routes

import (
	"net/http"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"

	"github.com/go-chi/chi"
)

type EventRouter struct {
	Handler events_handler.EventsHandlerI
}

func (er *EventRouter) Routes() http.Handler {
	r := chi.NewRouter()
	return r
}
