package user_handler

import (
	"context"
	"wiselink/internal/data/infrastructure/user_repository"
	"wiselink/pkg/Domain/user"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repository user_repository.UserRepositoryI
}
type UserHandlerI interface {
	GetByEmail(ctx context.Context, email string) (int, user.User)
	UserRegistration(ctx context.Context, u user.User) (int, user.User)
	DeleteUser(ctx context.Context, email string) int
}

func (uh *UserHandler) GetByEmail(ctx context.Context, email string) (int, user.User) {
	return uh.Repository.GetByEmail(ctx, email)
}

func (uh *UserHandler) UserRegistration(ctx context.Context, u user.User) (int, user.User) {
	id := uh.Repository.FindLastId(ctx)
	if id > -1 {
		u.Id = id + 1
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(u.TemporaryPassword+u.Email), bcrypt.MaxCost)
	if err != nil {
		return user.InternalError, u
	}
	u.AccessToken = string(pass)
	u.TemporaryPassword = ""
	status := uh.Repository.CreateUser(ctx, u)
	return status, u
}

func (uh *UserHandler) DeleteUser(ctx context.Context, email string) int {
	return uh.Repository.DeleteUser(ctx, email)
}
