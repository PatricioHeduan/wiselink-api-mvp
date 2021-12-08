package events_handler

import (
	"wiselink-api/data/internal/infrastructure/repositories/events_repository"
)

type EventsHandler struct {
	Repository events_repository.EventRepositoryI
}
type EventsHandlerI interface {
}
