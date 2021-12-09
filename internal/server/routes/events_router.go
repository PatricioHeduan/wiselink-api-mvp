package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wiselink/pkg/Domain/events"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"

	"github.com/go-chi/chi"
)

type EventRouter struct {
	Handler events_handler.EventsHandlerI
}

func (er *EventRouter) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e events.Event
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400:Bad Request"))
		return
	}
	defer r.Body.Close()
	//Todo: check if the user is an administrator
	createdEvent, status := er.Handler.CreateEvent(ctx, e)
	if status == events.Success {
		parsedEvent, err := json.Marshal(createdEvent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Servcer Error"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(parsedEvent))
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Servcer Error"))
		return
	}
}
func (er *EventRouter) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var e events.Event
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	defer r.Body.Close()
	//Todo: check if the user is an administrator
	switch er.Handler.UpdateEvent(ctx, e) {
	case events.Success:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: OK"))
		return
	case events.NotFound:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("404: Not Found"))
		return
	case events.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Servcer Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Servcer Error"))
		return
	}
}
func (er *EventRouter) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
	}
	defer r.Body.Close()
	//Todo: check if the user is an administrator
	switch er.Handler.DeleteEvent(ctx, id) {
	case events.Success:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: Success"))
		return
	case events.NotFound:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("404: Not Found"))
		return
	case events.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Servcer Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Servcer Error"))
		return
	}
}

func (er *EventRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/creteEvent", er.CreateEvent)
	r.Put("/updateEvent", er.UpdateEvent)
	r.Delete("/deleteEvent", er.DeleteEvent)

	return r
}
