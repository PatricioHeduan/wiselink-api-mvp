package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	events_repository "wiselink/internal/data/infrastructure/events_repository"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/user"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"
	user_handler "wiselink/pkg/Use_Cases/Handlers/user_handlers"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	Handler user_handler.UserHandlerI
}

var eh = &events_handler.EventsHandler{
	Repository: &events_repository.EventsRepository{
		Client: newClient,
	},
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
		case user.NotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404: Not Found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	token := r.Header.Get("Authorization")
	tokenStatus := ur.Handler.VerifyAdminExistance(ctx, token)
	switch tokenStatus {
	case user.Success:
		switch ur.Handler.DeleteUser(ctx, id) {
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
	case user.NotFound:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Unauthorized"))
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
	status, userFound := ur.Handler.GetUserById(ctx, u.Id)
	switch status {
	case user.Success:
		if userFound.Email == u.Email {
			switch ur.Handler.UpdateUser(ctx, u, userFound) {
			case user.Success:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200: OK"))
				return
			case user.NotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404: Not Found"))
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401: Email not matching"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
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
	token := r.Header.Get("Authorization")
	tokenStatus := ur.Handler.VerifyAdminExistance(ctx, token)
	switch tokenStatus {
	case user.Success:
		adminStatus, _ := ur.Handler.GetAdminByEmail(ctx, email)
		switch adminStatus {
		case user.NotFound:
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
		case user.Success:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401: Already Exists"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Unauthorized"))
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
	token := r.Header.Get("Authorization")
	tokenStatus := ur.Handler.VerifyAdminExistance(ctx, token)
	switch tokenStatus {
	case user.Success:
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
	case user.NotFound:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Unauthorized"))
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	defer r.Body.Close()
	userStatus, userFound := ur.Handler.GetUserById(ctx, userId)
	switch userStatus {
	case user.Success:
		eventStatus, eventFound := eh.GetEventById(ctx, eventId)
		switch eventStatus {
		case events.Success:
			status := ur.Handler.UserInscription(ctx, userFound, eventFound)
			switch status {
			case user.Success:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200: OK"))
				return
			case user.NotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404: User Not Found"))
				return
			case user.EventNotPublished:
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401: Unauthorized event not published"))
				return
			case user.CantEnroll:
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409: Conflict: Cant Enroll"))
				return
			case user.AlreadyInscripted:
				w.WriteHeader(http.StatusAlreadyReported)
				w.Write([]byte("208: AlreadyInscripted"))
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
				return
			}
		case events.NotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404: Event Not Found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) UserUnsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	defer r.Body.Close()
	userStatus, userFound := ur.Handler.GetUserById(ctx, userId)
	switch userStatus {
	case user.Success:
		eventStatus, eventFound := eh.GetEventById(ctx, eventId)
		switch eventStatus {
		case events.Success:
			status := ur.Handler.UserUnsubscribe(ctx, userFound, eventFound)
			switch status {
			case user.Success:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200: OK"))
				return
			case user.NotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404: User Not Found"))
				return
			case user.NotSuscripted:
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("403: Not Suscripted"))
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
				return
			}
		case events.NotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404: Event Not Found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	status, userFound := ur.Handler.GetByEmail(ctx, u.Email)
	switch status {
	case user.Success:
		loginStatus := ur.Handler.LoginUser(ctx, u, userFound.AccessToken)
		switch loginStatus {
		case user.Success:
			parsedUser, err := json.Marshal(userFound)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: Internal Server Error"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(parsedUser)
			return
		case user.IncorectPassword:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401: Incorrect Password"))
			return
		case user.NotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404: Not Found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404:Not Found"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) GetInscriptedEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	filter := r.URL.Query().Get("filter")
	if filter != "activo" && filter != "completado" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		return
	}
	userStatus, userFound := ur.Handler.GetUserById(ctx, id)
	switch userStatus {
	case user.Success:
		userInscriptions := ur.Handler.GetInscriptedEvents(ctx, filter, userFound, eh)
		inscriptionsParsed, err := json.Marshal(userInscriptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(inscriptionsParsed)
		return
	case user.NotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User Not Found"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
}

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/registUser", ur.UserRegistration)

	r.Delete("/deleteUser", ur.DeleteUser)
	r.Delete("/userUnsubscribe", ur.UserUnsubscribe)

	r.Put("/updateUser", ur.UpdateUser)
	r.Put("/userToAdmin", ur.UserToAdmin)
	r.Put("/adminToUser", ur.AdminToUser)
	r.Put("/userInscription", ur.UserInscription)

	r.Get("/login", ur.LoginUser)
	r.Get("/inscriptedEvents", ur.GetInscriptedEvents)

	return r
}
