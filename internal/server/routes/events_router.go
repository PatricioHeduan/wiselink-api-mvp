package routes

import (
	"net/http"
	events_handler "wiselink-api/pkg/Use_Cases/Handlers/events_handlers"

	"github.com/go-chi/chi"
)

type EventRouter struct {
	Handler events_handler.HandlerI
}

func (er *EventRouter) Routes() http.Handler {
	r := chi.NewRouter()
	return r
}
