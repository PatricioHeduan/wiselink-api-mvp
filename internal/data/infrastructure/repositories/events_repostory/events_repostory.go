package events_repository

import "go.mongodb.org/mongo-driver/mongo"

type EventRepository struct {
	Client *mongo.Client
}
type EventRepositoryI interface {
}
