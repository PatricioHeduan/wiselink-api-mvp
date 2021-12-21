package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wiselink/internal/data/infrastructure/user_repository"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/filters"
	"wiselink/pkg/Domain/user"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"
	user_handler "wiselink/pkg/Use_Cases/Handlers/user_handlers"

	"github.com/go-chi/chi"
)

type EventRouter struct {
	Handler events_handler.EventsHandlerI
}

var uh = &user_handler.UserHandler{
	Repository: &user_repository.UserRepository{
		Client: newClient,
	},
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
	//Todo: check if the user is an administrator
	createdEvent, status := er.Handler.CreateEvent(ctx, e)
	if status == events.Success {
		parsedEvent, err := json.Marshal(createdEvent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(parsedEvent))
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
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
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
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
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (er *EventRouter) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var f filters.Filter
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
	}
	defer r.Body.Close()
	token := r.Header.Get("Authorization")
	adminStatus := uh.VerifyAdminExistance(ctx, token)
	switch adminStatus {
	case user.Success, user.NotFound:
		var admin bool
		if adminStatus == user.Success {
			admin = true
		} else {
			admin = false
		}
		status, eventSlice := er.Handler.GetEvents(ctx, admin, f)
		switch status {
		case events.Success:
			parsedEvents, err := json.Marshal(eventSlice)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(parsedEvents)
			return
		case events.NotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404: Not Found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (er *EventRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/creteEvent", er.CreateEvent)
	r.Put("/updateEvent", er.UpdateEvent)
	r.Delete("/deleteEvent", er.DeleteEvent)
	r.Get("/", er.GetEvents)

	return r
}
