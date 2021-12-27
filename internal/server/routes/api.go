package routes

import (
	"context"
	"net/http"
	events_repository "wiselink/internal/data/infrastructure/events_repository"
	"wiselink/internal/data/infrastructure/user_repository"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"
	user_handler "wiselink/pkg/Use_Cases/Handlers/user_handlers"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientOptions = options.Client().ApplyURI("mongodb+srv://admin:12345@wiselink-mvp.u0dgh.mongodb.net/wlMVP?retryWrites=true&w=majority")
	newClient, _  = mongo.Connect(context.TODO(), clientOptions)
)

func New() http.Handler {
	r := chi.NewRouter()
	er := &EventRouter{
		Handler: &events_handler.EventsHandler{
			Repository: &events_repository.EventsRepository{
				Client: newClient,
			},
		},
	}

	ur := &UserRouter{
		Handler: &user_handler.UserHandler{
			Repository: &user_repository.UserRepository{
				Client: newClient,
			},
		},
	}

	r.Mount("/events", er.Routes())
	r.Mount("/users", ur.Routes())

	return r
}
