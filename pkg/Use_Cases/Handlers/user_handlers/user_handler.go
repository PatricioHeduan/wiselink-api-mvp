package user_handler

import "wiselink/internal/data/infrastructure/user_repository"

type UserHandler struct {
	Repository user_repository.UserRepositoryI
}
type UserHandlerI interface {
}
