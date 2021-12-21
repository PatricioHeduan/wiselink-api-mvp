package events_handler

import (
	"context"
	events_repository "wiselink/internal/data/infrastructure/events_repository"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/filters"
	helpers "wiselink/pkg/Use_Cases/Helpers"
)

type EventsHandler struct {
	Repository events_repository.EventsRepositoryI
}
type EventsHandlerI interface {
	CreateEvent(ctx context.Context, e events.Event) (events.Event, int)
	UpdateEvent(ctx context.Context, e events.Event) int
	DeleteEvent(ctx context.Context, id int) int
	GetEvents(ctx context.Context, admin bool, filter filters.Filter) (int, []events.Event)
}

func (eh *EventsHandler) CreateEvent(ctx context.Context, e events.Event) (events.Event, int) {
	id := eh.Repository.FindLastId(ctx)
	if id > -1 {
		e.Id = id + 1
		return e, eh.Repository.CreateEvent(ctx, e)
	}
	return e, events.InternalError
}

func (eh *EventsHandler) UpdateEvent(ctx context.Context, e events.Event) int {
	return eh.Repository.UpdateEvent(ctx, e)
}

func (eh *EventsHandler) DeleteEvent(ctx context.Context, id int) int {
	return eh.Repository.DeleteEvent(ctx, id)
}

func (eh *EventsHandler) GetEvents(ctx context.Context, admin bool, filter filters.Filter) (int, []events.Event) {
	status, eventSlice := eh.Repository.GetEvents(ctx)
	if status != events.Success {
		return status, eventSlice
	}
	eventsReturn := []events.Event{}
	for _, e := range eventSlice {
		if helpers.Filtered(e, filter) {
			if !admin {
				if e.Status {
					eventsReturn = append(eventsReturn, e)
				}
			} else {
				eventsReturn = append(eventsReturn, e)
			}
		}
	}
	return events.Success, eventsReturn
}
