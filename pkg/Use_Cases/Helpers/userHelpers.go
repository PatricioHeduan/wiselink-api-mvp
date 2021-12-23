package helpers

import (
	"context"
	"wiselink/internal/data/infrastructure/user_repository"
)

func GetUserLastId(ctx context.Context, repository user_repository.UserRepositoryI) int {
	id := repository.FindUserLastId(ctx)
	if id == -1 {
		return id
	}
	return id + 1
}

func GetAdminLastId(ctx context.Context, repository user_repository.UserRepositoryI) int {
	id := repository.GetLastAdminId(ctx)
	if id == -1 {
		return id
	}
	return id + 1
}
