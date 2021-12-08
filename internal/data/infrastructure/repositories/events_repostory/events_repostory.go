package events_repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Client *mongo.Client
}
type RepositoryI interface {
}
