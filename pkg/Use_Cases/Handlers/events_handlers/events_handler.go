package events_handler

import (
	events_repository "wiselink/internal/data/infrastructure/events_repository"
)

type EventsHandler struct {
	Repository events_repository.EventsRepositoryI
}
type EventsHandlerI interface {
}
