package user_handler

import (
	"context"
	"time"
	"wiselink/internal/data/infrastructure/user_repository"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/user"
	helpers "wiselink/pkg/Use_Cases/Helpers"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repository user_repository.UserRepositoryI
}
type UserHandlerI interface {
	GetByEmail(ctx context.Context, email string) (int, user.User)
	UserRegistration(ctx context.Context, u user.User) (int, user.User)
	DeleteUser(ctx context.Context, email string) int
	UpdateUser(ctx context.Context, u user.User) int
	UserToAdmin(ctx context.Context, u user.User) int
	AdminToUser(ctx context.Context, a user.Admin) int
	GetAdminByEmail(ctx context.Context, email string) (int, user.Admin)
	VerifyAdminExistance(ctx context.Context, accessToken string) int
	UserInscription(ctx context.Context, u user.User, e events.Event) int
	GetUserById(ctx context.Context, id int) (int, user.User)
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
	pass, err := bcrypt.GenerateFromPassword([]byte(u.TemporaryPassword+u.Email), bcrypt.MaxCost)
	if err != nil {
		return user.InternalError, u
	}
	u.Id = lastId
	u.AccessToken = string(pass)
	u.TemporaryPassword = ""
	status := uh.Repository.CreateUser(ctx, u)
	return status, u
}

func (uh *UserHandler) DeleteUser(ctx context.Context, email string) int {
	return uh.Repository.DeleteUser(ctx, email)
}

func (uh *UserHandler) UpdateUser(ctx context.Context, u user.User) int {
	return uh.Repository.UpdateUser(ctx, u)
}

func (uh *UserHandler) UserToAdmin(ctx context.Context, u user.User) int {
	lastId := helpers.GetAdminLastId(ctx, uh.Repository)
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
