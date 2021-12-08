package events_handler

import (
	"context"
	events_repository "wiselink/internal/data/infrastructure/events_repository"
	"wiselink/pkg/Domain/events"
)

type EventsHandler struct {
	Repository events_repository.EventsRepositoryI
}
type EventsHandlerI interface {
	CreateEvent(ctx context.Context, e events.Event) (events.Event, int)
}

func (eh *EventsHandler) CreateEvent(ctx context.Context, e events.Event) (events.Event, int) {
	id := eh.Repository.FindLastId(ctx)
	if id > -1 {
		e.Id = id + 1
		return e, eh.Repository.CreateEvent(ctx, e)
	}
	return e, events.InternalError
}
