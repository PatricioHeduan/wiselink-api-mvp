package routes

import (
	"context"
	"net/http"
	events_repository "wiselink-api/internal/data/repositories/events_repositories"
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
	er := EventRouter{
		Handler: &events_handler.Handler{
			Repository: &events_repository.Repository{
				Client: &newClient,
			},
		},
	}

	return r
}
