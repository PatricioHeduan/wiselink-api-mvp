package routes

import (
	"net/http"
	user_handler "wiselink/pkg/Use_Cases/Handlers/user_handlers"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	Handler user_handler.UserHandlerI
}

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()
	return r
}
