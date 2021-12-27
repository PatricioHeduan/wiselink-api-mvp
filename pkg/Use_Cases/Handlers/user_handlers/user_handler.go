package user_handler

import (
	"context"
	"time"
	"wiselink/internal/data/infrastructure/user_repository"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/user"
	events_handler "wiselink/pkg/Use_Cases/Handlers/events_handlers"
	helpers "wiselink/pkg/Use_Cases/Helpers"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repository user_repository.UserRepositoryI
}

type UserHandlerI interface {
	GetByEmail(ctx context.Context, email string) (int, user.User)
	UserRegistration(ctx context.Context, u user.User) (int, user.User)
	DeleteUser(ctx context.Context, id int) int
	UpdateUser(ctx context.Context, u, foundUser user.User) int
	UserToAdmin(ctx context.Context, u user.User) int
	AdminToUser(ctx context.Context, a user.Admin) int
	GetAdminByEmail(ctx context.Context, email string) (int, user.Admin)
	VerifyAdminExistance(ctx context.Context, accessToken string) int
	UserInscription(ctx context.Context, u user.User, e events.Event) int
	GetUserById(ctx context.Context, id int) (int, user.User)
	LoginUser(ctx context.Context, u user.User, token string) int
	UserUnsubscribe(ctx context.Context, u user.User, e events.Event) int
	GetInscriptedEvents(ctx context.Context, filter string, u user.User, e events_handler.EventsHandlerI) []events.Event
}

const (
	dateAndHourFormat = "02-01-2006 at 15:04"
)

func (uh *UserHandler) GetByEmail(ctx context.Context, email string) (int, user.User) {
	return uh.Repository.GetByEmail(ctx, email)
}

func (uh *UserHandler) UserRegistration(ctx context.Context, u user.User) (int, user.User) {
	lastId := helpers.GetUserLastId(ctx, uh.Repository)
	if lastId == -1 {
		return user.InternalError, u
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(u.TemporaryPassword+u.Email), bcrypt.DefaultCost)
	if err != nil {
		return user.InternalError, u
	}
	u.Id = lastId
	u.AccessToken = string(pass)
	u.TemporaryPassword = ""
	status := uh.Repository.CreateUser(ctx, u)
	return status, u
}

func (uh *UserHandler) DeleteUser(ctx context.Context, id int) int {
	return uh.Repository.DeleteUser(ctx, id)
}

func (uh *UserHandler) UpdateUser(ctx context.Context, u, foundUser user.User) int {
	token := foundUser.AccessToken
	if u.TemporaryPassword != "" {
		byteToken, err := bcrypt.GenerateFromPassword([]byte(foundUser.TemporaryPassword+u.Email), bcrypt.DefaultCost)
		if err != nil {
			return user.InternalError
		}
		token = string(byteToken)
	}
	return uh.Repository.UpdateUser(ctx, u, token)
}

func (uh *UserHandler) UserToAdmin(ctx context.Context, u user.User) int {
	lastId := helpers.GetAdminLastId(ctx, uh.Repository)
	if lastId == -1 {
		return user.InternalError
	}
	a := user.Admin{
		Id:          lastId,
		Email:       u.Email,
		AccessToken: u.AccessToken,
	}
	return uh.Repository.AddAdmin(ctx, a)
}

func (uh *UserHandler) AdminToUser(ctx context.Context, a user.Admin) int {
	return uh.Repository.DeleteAdmin(ctx, a)
}

func (uh *UserHandler) GetAdminByEmail(ctx context.Context, email string) (int, user.Admin) {
	return uh.Repository.GetAdminByEmail(ctx, email)
}

func (uh *UserHandler) VerifyAdminExistance(ctx context.Context, accessToken string) int {
	return uh.Repository.VerifyAdminExistance(ctx, accessToken)
}

func (uh *UserHandler) UserInscription(ctx context.Context, u user.User, e events.Event) int {
	for _, id := range u.SuscriptedTo {
		if id == e.Id {
			return user.AlreadyInscripted
		}
	}
	eventDateTime, err := time.Parse(dateAndHourFormat, e.Date+" at "+e.Hour)
	if err != nil {
		return user.InternalError
	}
	if eventDateTime.Before(time.Now()) {
		return user.CantEnroll
	}
	if !e.Status {
		return user.EventNotPublished
	}
	u.SuscriptedTo = append(u.SuscriptedTo, e.Id)
	return uh.Repository.ModifyUserEvents(ctx, u)
}

func (uh *UserHandler) GetUserById(ctx context.Context, id int) (int, user.User) {
	return uh.Repository.GetUserById(ctx, id)
}

func (uh *UserHandler) LoginUser(ctx context.Context, u user.User, token string) int {
	err := bcrypt.CompareHashAndPassword([]byte(token), []byte(u.TemporaryPassword+u.Email))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return user.IncorectPassword
		}
		return user.InternalError
	}
	return user.Success
}

func (uh *UserHandler) UserUnsubscribe(ctx context.Context, u user.User, e events.Event) int {
	found := -1
	for i, id := range u.SuscriptedTo {
		if id == e.Id {
			found = i
			break
		}
	}
	if found == -1 {
		return user.NotSuscripted
	}
	u.SuscriptedTo = append(u.SuscriptedTo[:found], u.SuscriptedTo[found+1:]...)
	return uh.Repository.ModifyUserEvents(ctx, u)
}

func (uh *UserHandler) GetInscriptedEvents(ctx context.Context, filter string, u user.User, e events_handler.EventsHandlerI) []events.Event {
	var eventsInscripted []events.Event
	for _, eventId := range u.SuscriptedTo {
		_, event := e.GetEventById(ctx, eventId)
		eventDateTime, _ := time.Parse(dateAndHourFormat, event.Date+" at "+event.Hour)
		if event.Status {
			if filter == "activo" {
				if time.Now().Before(eventDateTime) {
					eventsInscripted = append(eventsInscripted, event)
				}
			} else {
				if eventDateTime.Before(time.Now()) {
					eventsInscripted = append(eventsInscripted, event)
				}
			}
		}
	}
	return eventsInscripted
}
