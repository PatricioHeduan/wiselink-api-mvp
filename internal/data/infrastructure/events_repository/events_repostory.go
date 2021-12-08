package events_repository

import "go.mongodb.org/mongo-driver/mongo"

type EventsRepository struct {
	Client *mongo.Client
}
type EventsRepositoryI interface {
}
