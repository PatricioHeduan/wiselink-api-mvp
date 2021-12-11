package routes

import (
	"encoding/json"
	"net/http"
	"wiselink/pkg/Domain/user"
	user_handler "wiselink/pkg/Use_Cases/Handlers/user_handlers"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	Handler user_handler.UserHandlerI
}

func (ur *UserRouter) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var u user.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	status, _ := ur.Handler.GetByEmail(ctx, u.Email)
	switch status {
	case user.NotFound:
		creationStatus, createdUser := ur.Handler.UserRegistration(ctx, u)
		switch creationStatus {
		case user.Success:
			parsedUser, err := json.Marshal(createdUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(parsedUser)
			return
		}
	case user.Success:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Already Exists"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()
	return r
}