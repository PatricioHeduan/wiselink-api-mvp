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

//constant to parse from a string to a time.Time to compare dates
const (
	dateAndHourFormat = "02-01-2006 at 15:04"
)

//method to get an user from email
func (uh *UserHandler) GetByEmail(ctx context.Context, email string) (int, user.User) {
	return uh.Repository.GetByEmail(ctx, email)
}

//Method to create a user or return to router if an error ocurrs
func (uh *UserHandler) UserRegistration(ctx context.Context, u user.User) (int, user.User) {
	lastId := helpers.GetUserLastId(ctx, uh.Repository)
	if lastId == -1 {
		return user.InternalError, u
	}
	//Generate a hash from TemporaryPassword and user Email to put it in acesstoken
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

//Method to delete a user or return to router if an error ocurrs
func (uh *UserHandler) DeleteUser(ctx context.Context, id int) int {
	return uh.Repository.DeleteUser(ctx, id)
}

//Method to update a user (and generate a new hash to accesstoken if password is changed) or return to router if an error ocurrs
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

//Method to promote a user to admin
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

//Method to demote a user to admin
func (uh *UserHandler) AdminToUser(ctx context.Context, a user.Admin) int {
	return uh.Repository.DeleteAdmin(ctx, a)
}

//method to get an user from email
func (uh *UserHandler) GetAdminByEmail(ctx context.Context, email string) (int, user.Admin) {
	return uh.Repository.GetAdminByEmail(ctx, email)
}

//method to get a status number if admin exist and return it to the router
func (uh *UserHandler) VerifyAdminExistance(ctx context.Context, accessToken string) int {
	return uh.Repository.VerifyAdminExistance(ctx, accessToken)
}

//method to inscribe a user to a single event
func (uh *UserHandler) UserInscription(ctx context.Context, u user.User, e events.Event) int {
	for _, id := range u.SuscriptedTo {
		if id == e.Id {
			return user.AlreadyInscripted
		}
	}
	//changin the time from a string to a time.Time to compare now time to event date for know if a user can or not inscribe to a single event
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

//method to get an user from id
func (uh *UserHandler) GetUserById(ctx context.Context, id int) (int, user.User) {
	return uh.Repository.GetUserById(ctx, id)
}

//method to get an user after logging correctly
func (uh *UserHandler) LoginUser(ctx context.Context, u user.User, token string) int {
	//Decode the "accessToken" hash from TemporaryPassword and user Email to know if temporarypassword is or not correct
	err := bcrypt.CompareHashAndPassword([]byte(token), []byte(u.TemporaryPassword+u.Email))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return user.IncorectPassword
		}
		return user.InternalError
	}
	return user.Success
}

//method to unsuscribe a user to a single event
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

//method get all events of a user is inscripted
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
