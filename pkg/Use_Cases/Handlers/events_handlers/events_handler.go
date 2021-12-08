package events_handler

import (
	events_repository "wiselink-api/internal/data/repositories/events_repositories"
)

type Handler struct {
	Repository events_repository.RepositoryI
}
type HandlerI interface {
}
