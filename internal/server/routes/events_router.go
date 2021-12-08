package routes

import events_handler "wiselink-api/pkg/Use_Cases/Handlers/events_handlers"

type EventRouter struct {
	Handler events_handler.HandlerI
}
