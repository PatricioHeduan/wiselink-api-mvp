package routes

import (
	"context"
	"net/http"
	"wiselink-api/data/internal/infrastructure/repositories/events_repository"
	events_handler "wiselink-api/pkg/Use_Cases/Handlers/events_handlers"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientOptions = options.Client().ApplyURI("STRING DE CONEXION")
	newClient, _  = mongo.Connect(context.TODO(), clientOptions)
)

func New() http.Handler {
	r := chi.NewRouter()
	er := &EventRouter{
		Handler: &events_handler.EventsHandler{
			Repository: &events_repository.EventRepository{
				Client: &newClient,
			},
		},
	}
	er.Routes()
	return r
}
