package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
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
				return
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

func (ur *UserRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.URL.Query().Get("email")
	defer r.Body.Close()
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	//Todo: check if the user is an administrator
	switch ur.Handler.DeleteUser(ctx, email) {
	case user.Success:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: OK"))
		return
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found"))
		return
	case user.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}
func (ur *UserRouter) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	//Todo: check if the user is an administrator
	switch ur.Handler.UpdateUser(ctx, u) {
	case user.Success:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: OK"))
		return
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found"))
		return
	case user.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}
func (ur *UserRouter) UserToAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.URL.Query().Get("email")
	defer r.Body.Close()
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	status, userToPromote := ur.Handler.GetByEmail(ctx, email)
	switch status {
	case user.Success:
		switch ur.Handler.UserToAdmin(ctx, userToPromote) {
		case user.Success:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200: Success"))
			return
		case user.InternalError:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found"))
		return
	case user.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) AdminToUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.URL.Query().Get("email")
	defer r.Body.Close()
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	status, adminToUser := ur.Handler.GetAdminByEmail(ctx, email)
	switch status {
	case user.Success:
		switch ur.Handler.AdminToUser(ctx, adminToUser) {
		case user.Success:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200: Success"))
			return
		case user.InternalError:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found"))
		return
	case user.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) UserInscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	//eventId := strconv.Atoi(r.URL.Query().Get("eventId"))
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	defer r.Body.Close()
	userStatus, _ := ur.Handler.GetById(ctx, userId)

	switch userStatus {

	}
}

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", ur.UserRegistration)
	r.Delete("/", ur.DeleteUser)
	r.Put("/", ur.UpdateUser)
	r.Put("/userToAdmin", ur.UserToAdmin)
	r.Put("/adminToUser", ur.AdminToUser)

	return r
}
